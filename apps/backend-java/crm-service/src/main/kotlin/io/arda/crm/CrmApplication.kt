package io.arda.crm

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class CrmApplication

func main(args: Array<String>) {
    runApplication<CrmApplication>(*args)
}
// Trigger build
