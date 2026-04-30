package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AdministrativeUnitFilter struct {
	ParentID string
	Level    string
	PageFilter
}

type AdministrativeUnit struct {
	ID            string
	Code          string
	Name          string
	FullName      string
	ShortName     string
	Level         string
	UnitType      string
	ParentID      string
	Path          string
	SortOrder     int
	Latitude      float64
	Longitude     float64
	Status        string
	EffectiveFrom string
	EffectiveTo   string
	Source        string
	MetadataJSON  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type AdministrativeUnitNode struct {
	Unit     *AdministrativeUnit
	Children []*AdministrativeUnitNode
}

type AdministrativeUnitSyncResult struct {
	ProvinceCount int
	WardCount     int
	EffectiveDate string
	Source        string
}

func (uc *MdmUsecase) ListAdministrativeUnits(ctx context.Context, filter AdministrativeUnitFilter) ([]*AdministrativeUnit, string, error) {
	normalizePageFilter(&filter.PageFilter)
	return uc.repo.ListAdministrativeUnits(ctx, filter)
}

func (uc *MdmUsecase) ListAdministrativeUnitTree(ctx context.Context, filter AdministrativeUnitFilter) ([]*AdministrativeUnitNode, error) {
	filter.PageSize = 100
	filter.PageToken = ""
	var units []*AdministrativeUnit
	for {
		pageUnits, next, err := uc.repo.ListAdministrativeUnits(ctx, filter)
		if err != nil {
			return nil, err
		}
		units = append(units, pageUnits...)
		if next == "" {
			break
		}
		filter.PageToken = next
	}
	return buildAdministrativeUnitTree(units), nil
}

func (uc *MdmUsecase) ListProvinces(ctx context.Context, filter AdministrativeUnitFilter) ([]*AdministrativeUnit, string, error) {
	filter.Level = "PROVINCE"
	normalizePageFilter(&filter.PageFilter)
	return uc.repo.ListAdministrativeUnits(ctx, filter)
}

func (uc *MdmUsecase) ListWards(ctx context.Context, provinceID string, filter PageFilter) ([]*AdministrativeUnit, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListAdministrativeUnits(ctx, AdministrativeUnitFilter{
		ParentID:   provinceID,
		Level:      "WARD",
		PageFilter: filter,
	})
}

func (uc *MdmUsecase) GetAdministrativeUnit(ctx context.Context, id string) (*AdministrativeUnit, error) {
	return uc.repo.GetAdministrativeUnit(ctx, id)
}

func (uc *MdmUsecase) CreateAdministrativeUnit(ctx context.Context, unit *AdministrativeUnit) (*AdministrativeUnit, error) {
	normalizeAdministrativeUnit(unit)
	return uc.repo.CreateAdministrativeUnit(ctx, unit)
}

func (uc *MdmUsecase) UpdateAdministrativeUnit(ctx context.Context, unit *AdministrativeUnit) (*AdministrativeUnit, error) {
	normalizeAdministrativeUnit(unit)
	return uc.repo.UpdateAdministrativeUnit(ctx, unit)
}

func (uc *MdmUsecase) DeleteAdministrativeUnit(ctx context.Context, id string) error {
	return uc.repo.DeleteAdministrativeUnit(ctx, id)
}

func (uc *MdmUsecase) SyncAdministrativeUnitsFromAddressKit(ctx context.Context) (*AdministrativeUnitSyncResult, error) {
	const effectiveDate = "latest"
	const source = "CASSO AddressKit"

	provinceResp, err := fetchAddressKit[struct {
		RequestID string               `json:"requestId"`
		Provinces []addressKitProvince `json:"provinces"`
	}](ctx, effectiveDate+"/provinces")
	if err != nil {
		return nil, err
	}
	communeResp, err := fetchAddressKit[struct {
		RequestID string              `json:"requestId"`
		Communes  []addressKitCommune `json:"communes"`
	}](ctx, effectiveDate+"/communes")
	if err != nil {
		return nil, err
	}
	if len(provinceResp.Provinces) == 0 || len(communeResp.Communes) == 0 {
		return nil, fmt.Errorf("%w: addresskit returned empty administrative data", ErrInvalidArgument)
	}

	units := make([]*AdministrativeUnit, 0, len(provinceResp.Provinces)+len(communeResp.Communes))
	provinceCodes := make(map[string]struct{}, len(provinceResp.Provinces))
	for i, province := range provinceResp.Provinces {
		provinceCodes[province.Code] = struct{}{}
		units = append(units, addressKitProvinceToUnit(province, i+1, provinceResp.RequestID, source))
	}
	for i, commune := range communeResp.Communes {
		if _, ok := provinceCodes[commune.ProvinceCode]; !ok {
			return nil, fmt.Errorf("%w: province %s for commune %s not found", ErrInvalidArgument, commune.ProvinceCode, commune.Code)
		}
		units = append(units, addressKitCommuneToUnit(commune, i+1, communeResp.RequestID, source))
	}
	if err := uc.repo.ReplaceAdministrativeUnits(ctx, units); err != nil {
		return nil, err
	}
	return &AdministrativeUnitSyncResult{
		ProvinceCount: len(provinceResp.Provinces),
		WardCount:     len(communeResp.Communes),
		EffectiveDate: effectiveDate,
		Source:        source,
	}, nil
}

func normalizeAdministrativeUnit(unit *AdministrativeUnit) {
	unit.Code = strings.TrimSpace(unit.Code)
	unit.Name = strings.TrimSpace(unit.Name)
	unit.Level = upperDefault(unit.Level, "WARD")
	unit.UnitType = upperDefault(unit.UnitType, "XA")
	unit.Status = upperDefault(unit.Status, "ACTIVE")
	unit.MetadataJSON = jsonDefault(unit.MetadataJSON)
}

type addressKitProvince struct {
	Code                string `json:"code"`
	Name                string `json:"name"`
	EnglishName         string `json:"englishName"`
	AdministrativeLevel string `json:"administrativeLevel"`
	Decree              string `json:"decree"`
}

type addressKitCommune struct {
	Code                string `json:"code"`
	Name                string `json:"name"`
	EnglishName         string `json:"englishName"`
	AdministrativeLevel string `json:"administrativeLevel"`
	ProvinceCode        string `json:"provinceCode"`
	ProvinceName        string `json:"provinceName"`
	Decree              string `json:"decree"`
}

func fetchAddressKit[T any](ctx context.Context, path string) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://production.cas.so/address-kit/"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	client := &http.Client{Timeout: 45 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: addresskit returned HTTP %d", ErrInvalidArgument, resp.StatusCode)
	}
	var out T
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func addressKitProvinceToUnit(in addressKitProvince, sortOrder int, requestID, source string) *AdministrativeUnit {
	fullName := strings.TrimSpace(in.Name)
	unitType := addressKitUnitType(in.AdministrativeLevel, fullName)
	return &AdministrativeUnit{
		Code:          strings.TrimSpace(in.Code),
		Name:          shortAdministrativeName(fullName),
		FullName:      fullName,
		ShortName:     shortAdministrativeName(fullName),
		Level:         "PROVINCE",
		UnitType:      unitType,
		Path:          fullName,
		SortOrder:     sortOrder,
		Status:        "ACTIVE",
		EffectiveFrom: "",
		Source:        source,
		MetadataJSON:  addressKitMetadata(requestID, in.EnglishName, in.AdministrativeLevel, in.Decree, "", ""),
	}
}

func addressKitCommuneToUnit(in addressKitCommune, sortOrder int, requestID, source string) *AdministrativeUnit {
	fullName := strings.TrimSpace(in.Name)
	provinceName := strings.TrimSpace(in.ProvinceName)
	return &AdministrativeUnit{
		Code:          strings.TrimSpace(in.Code),
		Name:          shortAdministrativeName(fullName),
		FullName:      fullName,
		ShortName:     shortAdministrativeName(fullName),
		Level:         "WARD",
		UnitType:      addressKitUnitType(in.AdministrativeLevel, fullName),
		ParentID:      strings.TrimSpace(in.ProvinceCode),
		Path:          strings.TrimSpace(fullName + ", " + provinceName),
		SortOrder:     sortOrder,
		Status:        "ACTIVE",
		EffectiveFrom: "",
		Source:        source,
		MetadataJSON:  addressKitMetadata(requestID, in.EnglishName, in.AdministrativeLevel, in.Decree, in.ProvinceCode, provinceName),
	}
}

func addressKitUnitType(level, name string) string {
	value := strings.ToLower(strings.TrimSpace(level + " " + name))
	switch {
	case strings.Contains(value, "thành phố"):
		return "THANH_PHO"
	case strings.Contains(value, "phường"):
		return "PHUONG"
	case strings.Contains(value, "xã") || strings.Contains(value, "xa"):
		return "XA"
	case strings.Contains(value, "đặc khu") || strings.Contains(value, "dac khu"):
		return "DAC_KHU"
	default:
		return "TINH"
	}
}

func shortAdministrativeName(value string) string {
	value = strings.TrimSpace(value)
	for _, prefix := range []string{"Thành phố ", "Tỉnh ", "Phường ", "Xã ", "Đặc khu "} {
		if strings.HasPrefix(value, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(value, prefix))
		}
	}
	return value
}

func addressKitMetadata(requestID, englishName, level, decree, provinceCode, provinceName string) string {
	metadata := map[string]string{
		"source":               "CASSO AddressKit",
		"source_url":           "https://production.cas.so/address-kit/latest",
		"request_id":           requestID,
		"english_name":         englishName,
		"administrative_level": level,
		"decree":               decree,
	}
	if provinceCode != "" {
		metadata["province_code"] = provinceCode
		metadata["province_name"] = provinceName
	}
	data, err := json.Marshal(metadata)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func buildAdministrativeUnitTree(units []*AdministrativeUnit) []*AdministrativeUnitNode {
	nodes := make(map[string]*AdministrativeUnitNode, len(units))
	var roots []*AdministrativeUnitNode
	for _, unit := range units {
		nodes[unit.ID] = &AdministrativeUnitNode{Unit: unit}
	}
	for _, unit := range units {
		node := nodes[unit.ID]
		if unit.ParentID != "" {
			if parent := nodes[unit.ParentID]; parent != nil {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}
	return roots
}

func buildAreaTree(areas []*Area) []*AreaNode {
	nodes := make(map[string]*AreaNode, len(areas))
	var roots []*AreaNode
	for _, area := range areas {
		nodes[area.ID] = &AreaNode{Area: area}
	}
	for _, area := range areas {
		node := nodes[area.ID]
		if area.ParentID != "" {
			if parent := nodes[area.ParentID]; parent != nil {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}
	return roots
}
