package com.arda.labs.accounting

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class AccountingApplication

func main(args: Array<String>) {
    runApplication<AccountingApplication>(*args)
}
