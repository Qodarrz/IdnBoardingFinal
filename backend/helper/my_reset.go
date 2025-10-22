package helpers

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendTokenForgotEmail(toEmail string, token string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	subject := "Token Reset Password Anda"

	body := fmt.Sprintf(
		`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="font-family: 'Segoe UI', Helvetica, Arial, sans-serif; margin: 0; padding: 0; color: #333; background-color: #f7f7f7;">
    <table align="center" width="100%%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 30px auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.05);">
        <tr>
            <td style="padding: 40px 30px 20px 30px; text-align: center;">
                <div style="margin-bottom: 25px;">
                    <img src="https://ifwlvrfebaqcu0a3.public.blob.vercel-storage.com/images.png" alt="Company Logo" style="max-width: 180px; height: auto;">
                </div>
                <h2 style="color: #2c3e50; margin-bottom: 25px; font-weight: 600;">Reset Password Anda</h2>
                <p style="font-size: 16px; line-height: 1.6; margin-bottom: 20px;">Anda telah meminta untuk mereset password akun Anda. Gunakan token berikut untuk melanjutkan proses reset password:</p>
                
                <div style="
                    display: inline-block;
                    padding: 15px 30px;
                    background-color: #f8f9fa;
                    color: #e74c3c;
                    text-decoration: none;
                    border-radius: 6px;
                    font-weight: 600;
                    font-size: 22px;
                    margin: 25px 0;
                    letter-spacing: 2px;
                    border: 1px dashed #e74c3c;">
                    %s
                </div>
                
                <p style="font-size: 14px; color: #7f8c8d; line-height: 1.6; margin-bottom: 0;">Token ini akan kadaluarsa dalam 1 jam. Jika Anda tidak meminta reset password, abaikan email ini dan pertimbangkan untuk mengubah password akun Anda.</p>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px 30px; background-color: #f8f9fa; border-radius: 0 0 8px 8px; text-align: center; font-size: 13px; color: #7f8c8d;">
                <p style="margin: 0;">© 2025 SMK Negeri 4 Kota Bogor. Semua hak dilindungi.</p>
                <p style="margin: 10px 0 0 0;">Jika Anda memiliki pertanyaan, hubungi <a href="mailto:reizentechid@gmail.com" style="color: #3498db; text-decoration: none;">reizentechid@gmail.com</a></p>
            </td>
        </tr>
    </table>
</body>
</html>`, token)

	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	displayName := "SMKN 4 Kota Bogor"
	msg += fmt.Sprintf("From: %s <%s>\nTo: %s\nSubject: %s\n\n%s", displayName, from, toEmail, subject, body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, []byte(msg))
	if err != nil {
		return fmt.Errorf("gagal mengirim email: %w", err)
	}

	return nil
}