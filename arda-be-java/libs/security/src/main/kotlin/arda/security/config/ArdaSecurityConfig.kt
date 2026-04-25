package arda.security.config

import arda.security.filter.ArdaSecurityFilter
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.security.config.annotation.web.reactive.EnableWebFluxSecurity
import org.springframework.security.config.web.server.SecurityWebFiltersOrder
import org.springframework.security.config.web.server.ServerHttpSecurity
import org.springframework.security.web.server.SecurityWebFilterChain

@Configuration
@EnableWebFluxSecurity
class ArdaSecurityConfig {

    @Bean
    fun springSecurityFilterChain(http: ServerHttpSecurity): SecurityWebFilterChain {
        return http
            .csrf { it.disable() }
            .formLogin { it.disable() }
            .httpBasic { it.disable() }
            .authorizeExchange {
                // By default, let the Gateway handle auth, but we can add service-level rules here
                it.anyExchange().permitAll()
            }
            .addFilterAt(ArdaSecurityFilter(), SecurityWebFiltersOrder.AUTHENTICATION)
            .build()
    }
}
