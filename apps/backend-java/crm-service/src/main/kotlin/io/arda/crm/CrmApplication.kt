package io.arda.crm

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
open class CrmApplication

fun main(args: Array<String>) {
    runApplication<CrmApplication>(*args)
}
// Force rebuild 1
