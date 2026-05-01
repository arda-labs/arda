package io.arda.crm;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.data.jpa.repository.config.EnableJpaAuditing;
import org.springframework.data.jpa.repository.config.EnableJpaRepositories;

@SpringBootApplication(scanBasePackages = {"io.arda.crm", "arda.messaging"})
@EnableJpaRepositories(basePackages = "io.arda.crm.domain")
@EnableJpaAuditing
public class CrmApplication {

    public static void main(String[] args) {
        // Đảm bảo JVM sử dụng múi giờ chuẩn để tránh lỗi "Asia/Saigon" trên Windows
        java.util.TimeZone.setDefault(java.util.TimeZone.getTimeZone("Asia/Ho_Chi_Minh"));
        System.setProperty("user.timezone", "Asia/Ho_Chi_Minh");

        SpringApplication.run(CrmApplication.class, args);
    }
}
