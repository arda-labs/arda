package arda.common.util;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.text.NumberFormat;
import java.util.Locale;

/**
 * Utility for Money handling, synchronized with Go's decimal handling approach.
 * We prefer using String or Long (minor units) for transport and BigDecimal for calculation.
 */
public class MoneyUtils {
    private static final int DEFAULT_SCALE = 2;
    private static final RoundingMode DEFAULT_ROUNDING = RoundingMode.HALF_UP;

    public static BigDecimal toBigDecimal(String amount) {
        if (amount == null || amount.isBlank()) return BigDecimal.ZERO;
        return new BigDecimal(amount).setScale(DEFAULT_SCALE, DEFAULT_ROUNDING);
    }

    public static String toString(BigDecimal amount) {
        if (amount == null) return "0.00";
        return amount.setScale(DEFAULT_SCALE, DEFAULT_ROUNDING).toPlainString();
    }

    public static String format(BigDecimal amount) {
        return format(amount, new Locale("vi", "VN"));
    }

    public static String format(BigDecimal amount, Locale locale) {
        NumberFormat formatter = NumberFormat.getCurrencyInstance(locale);
        return formatter.format(amount);
    }

    public static BigDecimal add(BigDecimal a, BigDecimal b) {
        return a.add(b).setScale(DEFAULT_SCALE, DEFAULT_ROUNDING);
    }

    public static BigDecimal subtract(BigDecimal a, BigDecimal b) {
        return a.subtract(b).setScale(DEFAULT_SCALE, DEFAULT_ROUNDING);
    }
}
