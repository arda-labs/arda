package arda.accounting.config

import arda.common.context.ArdaContext
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.data.domain.ReactiveAuditorAware
import org.springframework.data.r2dbc.config.EnableR2dbcAuditing
import reactor.core.publisher.Mono

@Configuration
@EnableR2dbcAuditing
class ArdaConfig {

    /**
     * Integrates ArdaContext with Spring Data R2DBC Auditing.
     * Automatically fills 'created_by' and 'updated_by' fields.
     */
    @Bean
    fun auditorAware(): ReactiveAuditorAware<String> {
        return ReactiveAuditorAware {
            ArdaContext.current()
                .map { it.userId ?: "system" }
        }
    }
}
