// DESIGN PATTERN: Strategy Pattern + Template Method Pattern
package services

import (
	"fmt"
	"net/smtp"
	"sender-service/config"
	"sender-service/models"
)

// EmailService - Handles email operations with configurable strategies
type EmailService struct {
	config *config.Config // Composition: HAS-A configuration
}

// NewEmailService - Factory method with dependency injection
func NewEmailService(config *config.Config) *EmailService {
	return &EmailService{config: config}
}

// SendTransferEmail - Sends email notification for point transfers
func (s *EmailService) SendTransferEmail(transfer *models.Transfer) error {
	// STRATEGY PATTERN: Different authentication strategies
	var auth smtp.Auth

	if s.config.Email.GmailAddress != "" && s.config.Email.GmailAppPass != "" {
		// Strategy 1: Authenticated SMTP with Gmail
		auth = smtp.PlainAuth("", s.config.Email.GmailAddress, s.config.Email.GmailAppPass, s.config.Email.SMTPHost)
		fmt.Println("Using SMTP authentication")
	} else {
		// Strategy 2: Unauthenticated SMTP (for testing/development)
		fmt.Println("Warning: No SMTP credentials provided, attempting without authentication")
		auth = nil
	}

	// FRONTEND INTEGRATION: Generate claim URL with hash routing for SPA
	claimURL := fmt.Sprintf("%s/#/claim/%s", s.config.Frontend.URL, transfer.Token)

	subject := "You've Received Virtual Points!"

	//  TEMPLATE METHOD PATTERN: HTML email template
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            max-width: 600px; 
            margin: 0 auto; 
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            background: white;
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .header { 
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); 
            color: white; 
            padding: 30px; 
            text-align: center; 
        }
        .content { 
            padding: 30px; 
        }
        .button { 
            display: inline-block; 
            padding: 15px 30px; 
            background: #667eea; 
            color: white; 
            text-decoration: none; 
            border-radius: 5px; 
            margin: 20px 0; 
            font-size: 16px;
            font-weight: bold;
        }
        .points { 
            font-size: 24px; 
            font-weight: bold; 
            color: #667eea; 
        }
        .footer { 
            text-align: center; 
            padding: 20px; 
            color: #666; 
            font-size: 14px;
            background: #f9f9f9;
            border-top: 1px solid #eee;
        }
        .info-box {
            background: #fff3cd;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            border-left: 4px solid #ffc107;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1> You've Received Virtual Points!</h1>
        </div>
        <div class="content">
            <p>Hello <strong>%s</strong>,</p>
            <p>Great news! You have received <span class="points">%d virtual points</span> from <strong>%s</strong>.</p>
            
            <div style="text-align: center;">
                <a href="%s" class="button">Claim Your Points Now</a>
            </div>
            
            <div class="info-box">
                <p><strong> Important:</strong> This link will expire in 24 hours.</p>
                <p>If you don't have an account yet, you'll be able to create one after clicking the link.</p>
            </div>
            
            <p><strong>Email:</strong> Make sure to use <strong>%s</strong> when creating your account.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br><strong>Virtual Points Team</strong></p>
            <p style="font-size: 12px; color: #999;">This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
    `, transfer.ReceiverName, transfer.Points, transfer.SenderEmail, claimURL, transfer.ReceiverEmail)

	// EMAIL HEADERS: Professional email formatting
	headers := make(map[string]string)
	headers["From"] = s.config.Email.From
	headers["To"] = transfer.ReceiverEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""
	headers["X-Priority"] = "1"
	headers["Importance"] = "high"

	// MESSAGE CONSTRUCTION: Build RFC-compliant email
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// EMAIL DELIVERY: Send via SMTP
	err := smtp.SendMail(
		s.config.Email.SMTPHost+":"+s.config.Email.SMTPPort,
		auth,
		s.config.Email.From,
		[]string{transfer.ReceiverEmail},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("failed to send email to %s: %v", transfer.ReceiverEmail, err)
	}

	fmt.Printf(" Email sent successfully to: %s\n", transfer.ReceiverEmail)
	fmt.Printf("Claim URL: %s\n", claimURL)
	return nil
}
