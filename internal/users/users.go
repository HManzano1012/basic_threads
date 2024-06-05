package users

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"basicthreads/internal/database"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func LoginUser(email, password string) echo.Map {
	if len(email) == 0 || len(password) == 0 {
		response := echo.Map{
			"status":  "error",
			"code":    400,
			"message": "Email and password are required",
		}
		return response
	}

	authUser := database.AuthUser(email, password)

	if !authUser {
		response := echo.Map{
			"status":  "error",
			"code":    401,
			"message": "Invalid credentials",
		}

		return response
	}

	claims := &jwtCustomClaims{
		email,
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 4)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		response := echo.Map{
			"status":  "error",
			"code":    500,
			"message": "Internal server error",
			"error":   "internal_server_error",
		}

		return response
	}

	response := echo.Map{
		"status": "success",
		"code":   200,
		"token":  t,
	}

	return response
}

func RegisterUser(name, email, phone, password string) echo.Map {
	if len(name) == 0 || len(email) == 0 || len(phone) == 0 || len(password) == 0 {
		response := echo.Map{
			"status":  "error",
			"code":    400,
			"message": "Name, email,phone and password are required",
			"error":   "missing_fields",
		}

		return response
	}

	userExists := database.ValidateUserExists(email)
	if userExists {
		response := echo.Map{
			"status":  "error",
			"code":    400,
			"message": "User already exists",
			"error":   "user_exists",
		}
		return response
	}

	err := database.RegisterUser(name, email, phone, password)
	if err != nil {
		response := echo.Map{
			"status":  "error",
			"code":    500,
			"message": "Internal server error",
			"error":   "internal_server_error",
		}
		return response
	}

	sendMailRegister(email, name)
	response := echo.Map{
		"status":  "success",
		"code":    200,
		"message": "User registered successfully",
	}

	return response
}

func sendMailRegister(email, name string) {
	url := "https://api.brevo.com/v3/smtp/email"
	method := "POST"

	payload := strings.NewReader(`{  
   "sender":{  
      "name":"Sender Alex",
      "email":"senderalex@example.com"
   },
   "to":[  
      {  
         "email":"` + email + `",
         "name":"` + name + `"
      }
   ],
   "subject":"Hello world",
   "htmlContent":'<!doctype html>
<html>
  <body>
    <div
      style='background-color:#eff4f3;color:#242424;font-family:Charter, "Bitstream Charter", "Sitka Text", Cambria, serif;font-size:16px;font-weight:400;letter-spacing:0.15008px;line-height:1.5;margin:0;padding:32px 0;min-height:100%;width:100%'
    >
      <table
        align="center"
        width="100%"
        style="margin:0 auto;max-width:600px;background-color:#e9f4f3"
        role="presentation"
        cellspacing="0"
        cellpadding="0"
        border="0"
      >
      
        <tbody>
          <tr style="width:100%">
            <td>
              <div
                style="padding:0px 24px 0px 4px;background-color:#fcf8f8;text-align:center"
              >
                <a
                  href="https://es.shein.com"
                  style="text-decoration:none"
                  target="_blank"
                  ><img
                    alt="Threads"
                    src="img/threads.png"
                    width="200"
                    height="200"
                    style="width:200px;height:200px;outline:none;border:none;text-decoration:none;vertical-align:middle;display:inline-block;max-width:100%"
                /></a>
              </div>
              <div
                style="font-size:16px;font-weight:bold;text-align:center;padding:12px 24px 16px 24px"
              >
                Hola, ` + name + ` 👋,
              </div>
              <div
                style="color:#171717;background-color:#fefffc;font-size:16px;font-weight:bold;text-align:center;padding:12px 24px 12px 24px"
              >
                Gracias por registrarse en el sitio web THREADS.
              </div>
              <div
                style="font-size:15px;font-weight:normal;text-align:center;padding:16px 24px 16px 24px"
              >
                Al registrarte obtuvistes ciertos beneficios que el sitio web
                ofrece. DISFRUTA DE LA ROPA
              </div>
              <div style="padding:16px 0px 16px 0px">
                <hr
                  style="width:100%;border:none;border-top:1px solid #CCCCCC;margin:0"
                />
              </div>
              <div
                style="font-weight:bold;text-align:center;padding:0px 24px 0px 24px"
              >
                Cambiar clave de acceso:
              </div>
              <div style="text-align:center;padding:16px 24px 16px 24px">
                <a
                  href="https://www.usewaypoint.com"       
                  style="color:#0A0A0A;font-size:17px;font-weight:bold;background-color:#f8f5f5;border-radius:64px;display:inline-block;padding:4px 8px;text-decoration:none"
                  target="_blank"
                  ><span
                    ><!--[if mso
                      ]><i
                        style="letter-spacing: 8px;mso-font-width:-100%;mso-text-raise:12"
                        hidden
                        >&nbsp;</i
                      ><!
                    [endif]--></span
                  ><span>Click</span
                  ><span
                    ><!--[if mso
                      ]><i
                        style="letter-spacing: 8px;mso-font-width:-100%"
                        hidden
                        >&nbsp;</i
                      ><!
                    [endif]--></span
                  ></a
                >
              </div>
              <div
                style="font-size:13px;font-weight:bold;text-align:center;padding:16px 24px 16px 24px"
              >
                Para poder iniciar sesion, unicamente tienes que utilizar tu
                correo y la contraseña
              </div>
              <div style="text-align:center;padding:20px 24px 24px 24px">
                <a
                  href="https://www.usewaypoint.com"  
                  style="color:#0A0A0A;font-size:17px;font-weight:bold;background-color:#f4f8fa;border-radius:64px;display:block;padding:8px 12px;text-decoration:none"
                  target="_blank"
                  ><span
                    ><!--[if mso
                      ]><i
                        style="letter-spacing: 12px;mso-font-width:-100%;mso-text-raise:18"
                        hidden
                        >&nbsp;</i
                      ><!
                    [endif]--></span
                  ><span>INICIAR SESION </span
                  ><span
                    ><!--[if mso
                      ]><i
                        style="letter-spacing: 12px;mso-font-width:-100%"
                        hidden
                        >&nbsp;</i
                      ><!
                    [endif]--></span
                  ></a
                >
              </div>
              <div style="padding:16px 24px 40px 24px;text-align:center">
                <img
                  alt="Threads"
                  src="https://i.pinimg.com/564x/ca/40/2e/ca402e8a89e96630d45d3e38a0a6952c.jpg"
                  width="300"
                  height="300"
                  style="width:300px;height:300px;outline:none;border:none;text-decoration:none;vertical-align:middle;display:inline-block;max-width:100%"
                />
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </body>
</html>'
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add(
		"api-key",
		"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",
	)
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
