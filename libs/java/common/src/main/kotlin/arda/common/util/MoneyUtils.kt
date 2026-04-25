package arda.common.util

import java.math.BigDecimal
import java.math.RoundingMode
import java.text.NumberFormat
import java.util.*

/**
 * Utility for Money handling, synchronized with Go's decimal handling approach.
 * We prefer using String or Long (minor units) for transport and BigDecimal for calculation.
 */
object MoneyUtils {
    private const val DEFAULT_SCALE = 2
    private val DEFAULT_ROUNDING = RoundingMode.HALF_UP

    fun toBigDecimal(amount: String?): BigDecimal {
        if (amount.isNullOrBlank()) return BigDecimal.ZERO
        return BigDecimal(amount).setScale(DEFAULT_SCALE, DEFAULT_ROUNDING)
    }

    fun toString(amount: BigDecimal?): String {
        if (amount == null) return "0.00"
        return amount.setScale(DEFAULT_SCALE, DEFAULT_ROUNDING).toPlainString()
    }

    fun format(amount: BigDecimal, locale: Locale = Locale("vi", "VN")): String {
        val formatter = NumberFormat.getCurrencyInstance(locale)
        return formatter.format(amount)
    }

    fun add(a: BigDecimal, b: BigDecimal): BigDecimal = a.add(b).setScale(DEFAULT_SCALE, DEFAULT_ROUNDING)

    fun subtract(a: BigDecimal, b: BigDecimal): BigDecimal = a.subtract(b).setScale(DEFAULT_SCALE, DEFAULT_ROUNDING)
}
