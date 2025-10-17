package main

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"my_project/database"
	"my_project/server/handler"
	"net/http"
)

var db *pgx.Conn

var idClient int

func main() {
	db = database.Connect()
	defer database.Close(db)

	h := handler.NewHandler(db)

	http.HandleFunc("GET /person", h.GetPerson)
	http.HandleFunc("POST /register", h.CreateUser)
	http.HandleFunc("POST /login", h.CheckUser)
	http.HandleFunc("PUT /person/update", h.UpdatePerson)
	http.HandleFunc("DELETE /person/delete", h.DeletePerson)
	http.HandleFunc("GET /person/update", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Edit Profile</title>
</head>
<body>
    <h1>Edit Profile</h1>
    <form id="editForm">
        <!-- Скрытое поле для ID -->
        <input type="hidden" id="id">
        
        <input type="text" id="name" placeholder="Name">
        <input type="email" id="email" placeholder="Email">
        <input type="password" id="hash" placeholder="New Password">
        <button type="submit">Update</button>
    </form>
    
    <div id="message"></div>

    <script>
        document.addEventListener('DOMContentLoaded', function() {
            fetch('/person?name=Иван%20Петров')
                .then(response => response.json())
                .then(userData => {
                    document.getElementById('id').value = userData.id;
                    document.getElementById('name').value = userData.name || '';
                    document.getElementById('email').value = userData.email || '';
                })
                .catch(error => {
                    showMessage('Ошибка загрузки: ' + error.message, 'error');
                });
        });

        document.getElementById('editForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const formData = {
                id: parseInt(document.getElementById('id').value),
                name: document.getElementById('name').value,
                email: document.getElementById('email').value,
                hash: document.getElementById('hash').value
            };
            
            console.log('Отправляем данные:', formData);
            
            fetch('/person/update', {
                method: 'PUT',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(formData)
            })
            .then(response => response.json())
            .then(result => {
                showMessage('Данные обновлены!', 'success');
            })
            .catch(error => {
                showMessage('Ошибка: ' + error.message, 'error');
            });
        });

        function showMessage(text, type) {
            const messageDiv = document.getElementById('message');
            messageDiv.textContent = text;
            messageDiv.style.color = type === 'success' ? 'green' : 'red';
        }
    </script>
</body>
</html>
    `))
	})

	http.HandleFunc("GET /market", h.GetProducts)
	http.HandleFunc("GET /market/{id}", h.GetProduct)
	http.HandleFunc("GET /", h.Redirect)
	http.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
<h1>Registration Page</h1>
<form id="registerForm">
    <input type="text" id="name" name="name" placeholder="Name">
    <input type="email" id="email" name="email" placeholder="Email">
    <input type="password" id="password" name="password" placeholder="Password">
    <button type="submit">Register</button>
</form>

<p>Already have an account? <a href="/login">Login here</a></p>

<script>
document.getElementById('registerForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const formData = {
        name: document.getElementById('name').value,
        email: document.getElementById('email').value,
        hash: document.getElementById('password').value
    };
    
    console.log('1. Form data collected:', formData);
    
    fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData)
    })
    .then(response => {
        console.log('2. Response status:', response.status, response.statusText);
        console.log('3. Response redirected:', response.redirected);
        console.log('4. Response ok:', response.ok);
        
        if (response.ok) {
            console.log('5. Registration successful, redirecting...');
            window.location.href = '/login';
        } else {
            return response.text().then(text => {
                console.log('6. Error response text:', text);
                alert('Error: ' + text);
            });
        }
    })
    .catch(error => {
        console.log('7. Fetch error:', error);
        alert('Registration failed: ' + error);
    });
    
    console.log('8. Fetch request sent');
});
</script>
    `))
	})

	http.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 400px;
            margin: 100px auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .login-container {
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            text-align: center;
            color: #333;
            margin-bottom: 30px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            color: #555;
            font-weight: bold;
        }
        input[type="email"],
        input[type="password"] {
            width: 100%;
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 4px;
            box-sizing: border-box;
            font-size: 16px;
        }
        input[type="email"]:focus,
        input[type="password"]:focus {
            outline: none;
            border-color: #4CAF50;
        }
        button {
            width: 100%;
            padding: 12px;
            background-color: #4CAF50;
            color: white;
            border: none;
            border-radius: 4px;
            font-size: 16px;
            cursor: pointer;
        }
        button:hover {
            background-color: #45a049;
        }
        button:disabled {
            background-color: #cccccc;
            cursor: not-allowed;
        }
        .links {
            text-align: center;
            margin-top: 20px;
        }
        .links a {
            color: #4CAF50;
            text-decoration: none;
        }
        .links a:hover {
            text-decoration: underline;
        }
        .message {
            padding: 10px;
            border-radius: 4px;
            margin-bottom: 20px;
            text-align: center;
            display: none;
        }
        .error {
            background-color: #ffebee;
            color: #c62828;
            border: 1px solid #ffcdd2;
        }
        .success {
            background-color: #e8f5e8;
            color: #2e7d32;
            border: 1px solid #c8e6c9;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <h1>Login</h1>
        
        <div id="message"></div>
        
        <form id="loginForm">
            <div class="form-group">
                <label for="email">Email:</label>
                <input type="email" id="email" name="email" placeholder="Enter your email" required>
            </div>
            
            <div class="form-group">
                <label for="password">Password:</label>
                <input type="password" id="password" name="password" placeholder="Enter your password" required>
            </div>
            
            <button type="submit">Login</button>
        </form>
        
        <div class="links">
            <p>Don't have an account? <a href="/register">Register here</a></p>
        </div>
    </div>

    <script>
        document.getElementById('loginForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            // Показываем loading state
            const button = document.querySelector('button');
            const originalText = button.textContent;
            button.textContent = 'Logging in...';
            button.disabled = true;
            
            const formData = {
                email: document.getElementById('email').value,
                hash: document.getElementById('password').value
            };
            
            console.log('Sending JSON:', JSON.stringify(formData, null, 2));
            
            // Отправляем JSON
            fetch('/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                'Accept': 'application/json'
                },
                body: JSON.stringify(formData)
            })
            .then(response => {
                console.log('Response status:', response.status);
                console.log('Content-Type:', response.headers.get('Content-Type'));
                
                // Пытаемся получить JSON ответ
                return response.json().then(data => {
                    return {
                        status: response.status,
                        ok: response.ok,
                        data: data
                    };
                });
            })
            .then(({ status, ok, data }) => {
                console.log('Response data:', data);
                
                if (ok) {
                    showMessage('Login successful!', 'success');
                    // Если в ответе есть redirect URL, используем его
                    if (data.redirect) {
                        setTimeout(() => {
                            window.location.href = data.redirect;
                        }, 1000);
                    } else {
                        setTimeout(() => {
                            window.location.href = '/dashboard';
                        }, 1000);
                    }
                } else {
                    showMessage(data.error || data.message || 'Login failed', 'error');
                }
            })
            .catch(error => {
                console.error('Login error:', error);
                showMessage('Login failed: ' + error.message, 'error');
            })
            .finally(() => {
                // Восстанавливаем кнопку
                button.textContent = originalText;
                button.disabled = false;
            });
        });

        function showMessage(text, type) {
            const messageDiv = document.getElementById('message');
            messageDiv.textContent = text;
            messageDiv.className = 'message ' + type;
            messageDiv.style.display = 'block';
        }

        // Показываем сообщение из URL параметров (например, после регистрации)
        const urlParams = new URLSearchParams(window.location.search);
        const message = urlParams.get('message');
        if (message) {
            showMessage(message, 'success');
        }
    </script>
</body>
</html>
    `))
	})
	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		log.Println(err)
	}
}
