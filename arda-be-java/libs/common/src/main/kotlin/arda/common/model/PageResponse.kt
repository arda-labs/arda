package arda.common.model

data class PageResponse<T>(
    val items: List<T>,
    val metadata: PageMetadata
) {
    companion object {
        fun <T> of(items: List<T>, page: Int, size: Int, totalElements: Long): PageResponse<T> {
            val totalPages = if (size == 0) 0 else Math.ceil(totalElements.toDouble() / size).toInt()
            return PageResponse(
                items = items,
                metadata = PageMetadata(
                    page = page,
                    size = size,
                    totalElements = totalElements,
                    totalPages = totalPages,
                    hasNext = page < totalPages - 1,
                    hasPrevious = page > 0
                )
            )
        }
    }
}

data class PageMetadata(
    val page: Int,
    val size: Int,
    val totalElements: Long,
    val totalPages: Int,
    val hasNext: Boolean,
    val hasPrevious: Boolean
)
