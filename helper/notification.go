package helper

import (
	"github.com/go-gomail/gomail"
	"io"
	"os"
	"strconv"
)

func SendWelcomeEmail(doctorEmail, fullName, verificationToken string) error {
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	sender := smtpUsername
	recipient := doctorEmail
	subject := "Welcome to Prodia"
	verificationLink := "http://35.225.10.188:8080/verify?token=" + verificationToken
	emailBody := `
    <html>
    <head>
        <link href="https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css" rel="stylesheet">
        <style>
            body {
                font-family: 'Arial', sans-serif;
                background-color: #f5f5f5;
                margin: 0;
                padding: 0;
            }
            .container {
                max-width: 600px;
                margin: 0 auto;
                padding: 20px;
                background-color: #ffffff;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                border-radius: 5px;
            }
            h1 {
                text-align: center;
                color: #333;
            }
            .message {
                background-color: #f9f9f9;
                padding: 15px;
                border: 1px solid #ddd;
                border-radius: 5px;
            }
            p {
                font-size: 16px;
                margin-top: 10px;
                line-height: 1.6;
            }
            strong {
                font-weight: bold;
            }
            .footer {
                text-align: center;
                margin-top: 20px;
                color: #666;
            }
            .btn-verify-email {
                background-color: #1E90FF;
                color: #fff;
                padding: 10px 20px;
                border-radius: 5px;
                text-decoration: none;
                display: inline-block;
                margin: 20px auto;
            }
            .btn-verify-email:hover {
                background-color: #007BFF;
            }
            .logo {
                text-align: center;
                margin-top: 20px;
            }
            .logo img {
                width: 120px;
                height: 120px;
                border-radius: 50%;
                border: 3px solid #1E90FF;
                transition: transform 0.3s ease-in-out;
                margin: 0 auto;
                display: block;
            }
            .logo img:hover {
                transform: scale(1.1);
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="logo">
                <a href="https://ibb.co.com/txfvnNV"><img src="https://i.ibb.co.com/LJwcBqV/channels4-profile.jpg" alt="channels4-profile" border="0"></a>
            </div>
            <h1>Welcome to Prodia</h1>
            <div class="message">
                <p>Hello, <strong>` + fullName + `</strong>,</p>
                <p>Thank you for choosing Prodia. You're now part of our team!</p>
                <p>If you have any questions or need assistance, please don't hesitate to contact our support team.</p>
                <p><strong>Support Team:</strong> <a href="mailto:prodiaaimar@gmail.com">prodiaaimar@gmail.com</a></p>
                <a href="` + verificationLink + `" class="btn btn-verify-email">Verify Email</a>
            </div>
            <div class="footer">
                <p>&copy; 2024 Prodia. All rights reserved.</p>
            </div>
        </div>
    </body>
    </html>
    `

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendLoginNotification(doctorEmail string, name string) error {
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	sender := smtpUsername
	recipient := doctorEmail
	subject := "Successful Login Notification"
	emailBody := `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login Notification</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            background: linear-gradient(180deg, #007BFF, #00BFFF);
            color: #fff;
            margin: 0;
            padding: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            height: 100vh;
        }
        .container {
            max-width: 600px;
            width: 100%;
            background-color: #fff;
            box-shadow: 0 0 20px rgba(0, 0, 0, 0.2);
            border-radius: 10px;
            overflow: hidden;
            text-align: center;
            margin: 0 auto; /* Menempatkan container di tengah */
        }
        .header {
            background-color: #007BFF;
            color: #fff;
            padding: 20px;
            border-bottom: 1px solid #ddd;
        }
        h1 {
            margin: 0;
            color: #333;
            font-size: 28px;
        }
        .logo {
            text-align: center;
            margin-top: 20px;
        }
        .logo img {
            width: 120px;
            height: 120px;
            border-radius: 50%;
            border: 3px solid #007BFF;
            transition: transform 0.3s ease-in-out;
        }
        .logo img:hover {
            transform: scale(1.1);
        }
        .message {
            padding: 20px;
        }
        p {
            font-size: 18px;
            margin-top: 15px;
            color: #555;
            line-height: 1.5;
        }
        .footer {
            text-align: center;
            padding: 20px;
            color: #666;
            font-size: 14px;
            border-top: 1px solid #ddd;
        }
        a {
            text-decoration: none;
            color: #007BFF;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Login Successful</h1>
        </div>
       <div class="logo">
    <img src="https://ibb.co.com/txfvnNV" alt="Prodia Logo">
</div>

        <div class="message">
            <p>Hello, <strong>` + name + `</strong>,</p>
            <p>Your login was successful. If this wasn't you, please contact our support team immediately. Thank you.</p>
            <p><strong>Support Team:</strong> <a href="mailto:prodiaaimar@gmail.com">prodiaaimar@gmail.com</a></p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Prodia. All rights reserved. | <a href="https://hr-harmony.seculab.space" target="_blank">Prodia</a></p>
        </div>
    </div>
</body>
</html>




	`

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendMedicalRecordNotification(patientEmail, patientName, diagnosis, prescription, careSuggestion string) error {
	pdfBytes, err := GenerateMedicalRecordPDF(patientName, diagnosis, prescription, careSuggestion)
	if err != nil {
		return err
	}

	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	sender := smtpUsername
	recipient := patientEmail
	subject := "Your Medical Record from Prodia"
	emailBody := `
	<html>
	<head>
		<style>
			/* Styles for email body */
		</style>
	</head>
	<body>
		<p>Hello, <strong>` + patientName + `</strong>,</p>
		<p>Please find attached your medical record from Prodia.</p>
		<p>If you have any questions or need assistance, please contact us at prodiaaimar@gmail.com.</p>
		<p>Regards,<br>Prodia Team</p>
	</body>
	</html>
	`

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody)
	m.Attach("medical_record.pdf", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(pdfBytes)
		return err
	}))

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
