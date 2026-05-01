package arda.common.model;

import java.util.List;

public record PageResponse<T>(
    List<T> items,
    PageMetadata metadata
) {
    public static <T> PageResponse<T> of(List<T> items, int page, int size, long totalElements) {
        int totalPages = (size == 0) ? 0 : (int) Math.ceil((double) totalElements / size);
        return new PageResponse<>(
            items,
            new PageMetadata(
                page,
                size,
                totalElements,
                totalPages,
                page < totalPages - 1,
                page > 0
            )
        );
    }

    public record PageMetadata(
        int page,
        int size,
        long totalElements,
        int totalPages,
        boolean hasNext,
        boolean hasPrevious
    ) {}
}
