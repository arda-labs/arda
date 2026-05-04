package data

import (
	"context"
	"fmt"
	"time"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/biz"
)

type eventRepo struct {
	data *Data
}

func NewEventRepo(data *Data) biz.EventRepo {
	return &eventRepo{data: data}
}

func (r *eventRepo) ListByInstance(ctx context.Context, instanceID string, pageSize int, pageToken string) ([]*biz.ProcessEvent, string, error) {
	params := pagination.Normalize(pageSize, pageToken)

	query := `SELECT id, process_instance_id, event_type, source, data, timestamp
	           FROM process_events
	           WHERE process_instance_id = $1
	           ORDER BY timestamp DESC
	           LIMIT $2 OFFSET $3`

	rows, err := r.data.DB(ctx).Pool.Query(ctx, query, instanceID, params.Limit, params.Offset)
	if err != nil {
		return nil, "", fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	var events []*biz.ProcessEvent
	for rows.Next() {
		e := &biz.ProcessEvent{}
		if err := rows.Scan(&e.ID, &e.ProcessInstanceID, &e.EventType, &e.Source, &e.Data, &e.Timestamp); err != nil {
			return nil, "", fmt.Errorf("scan event: %w", err)
		}
		events = append(events, e)
	}

	nextToken := ""
	if len(events) > 0 {
		nextToken = pagination.NextOffsetToken(len(events), params.Limit, params.Offset)
	}

	return events, nextToken, rows.Err()
}

func (r *eventRepo) Create(ctx context.Context, event *biz.ProcessEvent) (*biz.ProcessEvent, error) {
	query := `INSERT INTO process_events (process_instance_id, event_type, source, data, timestamp)
	          VALUES ($1, $2, $3, $4::jsonb, $5)
	          RETURNING id`

	created := &biz.ProcessEvent{}
	err := r.data.DB(ctx).Pool.QueryRow(ctx, query,
		event.ProcessInstanceID, event.EventType, event.Source, event.Data, time.Now()).Scan(&created.ID)
	if err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}
	created.ProcessInstanceID = event.ProcessInstanceID
	created.EventType = event.EventType
	created.Source = event.Source
	created.Data = event.Data
	created.Timestamp = time.Now()
	return created, nil
}

func (r *eventRepo) GetHeatmap(ctx context.Context, definitionID string) ([]*biz.HeatmapStep, error) {
	query := `
		SELECT
			data->>'elementId' AS element_id,
			COUNT(*) AS instance_count,
			AVG(EXTRACT(EPOCH FROM timestamp - lag.timestamp)) AS avg_duration_seconds
		FROM process_events pe
		JOIN process_instances pi ON pi.id = pe.process_instance_id
		LEFT JOIN LATERAL (
			SELECT timestamp FROM process_events pe2
			WHERE pe2.process_instance_id = pe.process_instance_id
			  AND pe2.event_type = pe.event_type
			  AND pe2.timestamp < pe.timestamp
			ORDER BY pe2.timestamp DESC
			LIMIT 1
		) lag ON true
		WHERE pi.process_definition_id = $1
		  AND pe.event_type IN ('ELEMENT_COMPLETED', 'ELEMENT_ACTIVATED')
		GROUP BY data->>'elementId'
		ORDER BY instance_count DESC`

	rows, err := r.data.DB(ctx).Pool.Query(ctx, query, definitionID)
	if err != nil {
		return nil, fmt.Errorf("query heatmap: %w", err)
	}
	defer rows.Close()

	var steps []*biz.HeatmapStep
	for rows.Next() {
		s := &biz.HeatmapStep{}
		if err := rows.Scan(&s.ElementID, &s.InstanceCount, &s.AvgDurationSeconds); err != nil {
			return nil, fmt.Errorf("scan heatmap step: %w", err)
		}
		// Determine severity
		switch {
		case s.InstanceCount > 100:
			s.Severity = "danger"
		case s.InstanceCount > 50:
			s.Severity = "warn"
		default:
			s.Severity = "success"
		}
		steps = append(steps, s)
	}
	return steps, rows.Err()
}
