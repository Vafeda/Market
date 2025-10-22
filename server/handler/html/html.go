package html

import (
	"fmt"
	"net/http"
)

const commonStyles = `
<style>
    * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
    }
    body {
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        min-height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 20px;
    }
    .container {
        background: white;
        padding: 40px;
        border-radius: 12px;
        box-shadow: 0 15px 35px rgba(0,0,0,0.1);
        width: 100%;
        max-width: 450px;
    }
    h1 {
        text-align: center;
        color: #333;
        margin-bottom: 30px;
        font-size: 28px;
        font-weight: 600;
    }
    .form-group {
        margin-bottom: 20px;
    }
    label {
        display: block;
        margin-bottom: 8px;
        color: #555;
        font-weight: 500;
        font-size: 14px;
    }
    input[type="text"],
    input[type="email"],
    input[type="password"],
    input[type="number"] {
        width: 100%;
        padding: 12px 16px;
        border: 2px solid #e1e5e9;
        border-radius: 8px;
        font-size: 16px;
        transition: all 0.3s ease;
        background-color: #f8f9fa;
    }
    input:focus {
        outline: none;
        border-color: #667eea;
        background-color: white;
        box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }
    button {
        width: 100%;
        padding: 14px;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        border: none;
        border-radius: 8px;
        font-size: 16px;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.3s ease;
        margin-top: 10px;
    }
    button:hover {
        transform: translateY(-2px);
        box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
    }
    button:disabled {
        opacity: 0.6;
        cursor: not-allowed;
        transform: none;
    }
    .links {
        text-align: center;
        margin-top: 25px;
        padding-top: 20px;
        border-top: 1px solid #e1e5e9;
    }
    .links a {
        color: #667eea;
        text-decoration: none;
        font-weight: 500;
        transition: color 0.3s ease;
    }
    .links a:hover {
        color: #764ba2;
        text-decoration: underline;
    }
    .message {
        padding: 12px;
        border-radius: 8px;
        margin-bottom: 20px;
        text-align: center;
        font-weight: 500;
        display: none;
    }
    .error {
        background-color: #ffeaea;
        color: #d63031;
        border: 1px solid #ffcccc;
    }
    .success {
        background-color: #e8f5e8;
        color: #2e7d32;
        border: 1px solid #c8e6c9;
    }
    .hidden {
        display: none;
    }
</style>
`

const commonScript = `
<script>
    function showMessage(text, type) {
        const messageDiv = document.getElementById('message');
        messageDiv.textContent = text;
        messageDiv.className = 'message ' + type;
        messageDiv.style.display = 'block';
        
        // Автоскрытие для успешных сообщений
        if (type === 'success') {
            setTimeout(() => {
                messageDiv.style.display = 'none';
            }, 5000);
        }
    }
    
    function setLoading(button, isLoading) {
        const originalText = button.textContent;
        if (isLoading) {
            button.innerHTML = '<span class="loading">⏳</span> Loading...';
            button.disabled = true;
        } else {
            button.textContent = originalText;
            button.disabled = false;
        }
    }
</script>
`

func UserUpdatePage(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Edit Profile - Market</title>
		%s
	</head>
	<body>
		<div class="container">
			<h1>Edit Profile</h1>
			
			<div id="message"></div>
			
			<form id="editForm">
				<input type="hidden" id="id">
				
				<div class="form-group">
					<label for="name">Full Name:</label>
					<input type="text" id="name" placeholder="Enter your name" required>
				</div>
				
				<div class="form-group">
					<label for="login">Login:</label>
					<input type="text" id="login" placeholder="Enter your login" required>
				</div>
				
				<div class="form-group">
					<label for="email">Email:</label>
					<input type="email" id="email" placeholder="Enter your email" required>
				</div>
				
				<div class="form-group">
					<label for="hash">New Password (optional):</label>
					<input type="password" id="hash" placeholder="Leave empty to keep current password">
				</div>
				
				<button type="submit">Update Profile</button>
			</form>
			
			<div class="links">
				<p><a href="/market">← Back to Market</a></p>
			</div>
		</div>
	
		%s
	
		<script>
			// Загружаем данные пользователя
			async function loadUserData() {
				try {
					console.log('Fetching from /user/1...');
					const response = await fetch('/user/1');
					console.log('Response status:', response.status);
					
					const text = await response.text();
					console.log('Response text (first 100 chars):', text.substring(0, 100));
					
					if (!response.ok) throw new Error('Failed to load user data');
					
					const userData = JSON.parse(text);
					console.log('User data:', userData);
					
					// ОТЛАДКА: проверяем элементы формы
					console.log('Form elements:');
					console.log('id element:', document.getElementById('id'));
					console.log('name element:', document.getElementById('name'));
					console.log('login element:', document.getElementById('login'));
					console.log('email element:', document.getElementById('email'));
					
					// Заполняем форму
					document.getElementById('id').value = userData.id;
					document.getElementById('name').value = userData.name || '';
					document.getElementById('login').value = userData.login || '';
					document.getElementById('email').value = userData.email || '';
					
					console.log('Form filled successfully');
					
				} catch (error) {
					console.error('Full error:', error);
					showMessage('Error loading user data: ' + error.message, 'error');
				}
			}
	
			document.addEventListener('DOMContentLoaded', loadUserData);
	
			document.getElementById('editForm').addEventListener('submit', async function(e) {
				e.preventDefault();
				
				const button = this.querySelector('button');
				setLoading(button, true);
				
				const formData = {
					id: parseInt(document.getElementById('id').value),
					name: document.getElementById('name').value,
					login: document.getElementById('login').value,
					email: document.getElementById('email').value
				};
				
				// Добавляем пароль только если он введен
				const password = document.getElementById('hash').value;
				if (password) {
					formData.hash = password;
				}
				
				try {
					const response = await fetch('/user/update', {
						method: 'PUT',
						headers: {'Content-Type': 'application/json'},
						body: JSON.stringify(formData)
					});
					
					if (response.ok) {
						showMessage('Profile updated successfully!', 'success');
					} else {
						const error = await response.text();
						showMessage(error || 'Update failed', 'error');
					}
				} catch (error) {
					showMessage('Update failed: ' + error.message, 'error');
				} finally {
					setLoading(button, false);
				}
			});
		</script>
	</body>
	</html>
    `, commonStyles, commonScript)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func RegistrationPage(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Register - Market</title>
		%s
	</head>
	<body>
		<div class="container">
			<h1>Create Account</h1>
			
			<div id="message"></div>
			
			<form id="registerForm">
				<div class="form-group">
					<label for="name">Full Name:</label>
					<input type="text" id="name" name="name" placeholder="Enter your full name" required>
				</div>
				
				<div class="form-group">
					<label for="login">Login:</label>
					<input type="text" id="login" name="login" placeholder="Choose a username" required>
				</div>
				
				<div class="form-group">
					<label for="email">Email:</label>
					<input type="email" id="email" name="email" placeholder="Enter your email" required>
				</div>
				
				<div class="form-group">
					<label for="password">Password:</label>
					<input type="password" id="password" name="password" placeholder="Create a password" required minlength="6">
				</div>
				
				<button type="submit">Create Account</button>
			</form>
			
			<div class="links">
				<p>Already have an account? <a href="/login">Sign in here</a></p>
			</div>
		</div>
	
		%s
	
		<script>
			document.getElementById('registerForm').addEventListener('submit', async function(e) {
				e.preventDefault();
				
				const button = this.querySelector('button');
				setLoading(button, true);
				
				const formData = {
					name: document.getElementById('name').value,
					login: document.getElementById('login').value,
					email: document.getElementById('email').value,
					hash: document.getElementById('password').value
				};
				
				try {
					const response = await fetch('/register', {
						method: 'POST',
						headers: {'Content-Type': 'application/json'},
						body: JSON.stringify(formData)
					});
					
					if (response.ok) {
						showMessage('Account created successfully! Redirecting...', 'success');
						setTimeout(() => window.location.href = '/login', 1500);
					} else {
						const error = await response.text();
						showMessage(error || 'Registration failed', 'error');
					}
				} catch (error) {
					showMessage('Registration failed: ' + error.message, 'error');
				} finally {
					setLoading(button, false);
				}
			});
		</script>
	</body>
	</html>
    `, commonStyles, commonScript)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - Market</title>
    %s
</head>
<body>
    <div class="container">
        <h1>Welcome Back</h1>
        
        <div id="message"></div>
        
        <form id="loginForm">
            <div class="form-group">
                <label for="login">Email or Login:</label>
                <input type="text" id="login" name="login" placeholder="Enter your email or login" required>
            </div>
            
            <div class="form-group">
                <label for="password">Password:</label>
                <input type="password" id="password" name="password" placeholder="Enter your password" required>
            </div>
            
            <button type="submit">Sign In</button>
        </form>
        
        <div class="links">
            <p>Don't have an account? <a href="/register">Create one here</a></p>
        </div>
    </div>

    %s

    <script>
        document.getElementById('loginForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const button = this.querySelector('button');
            setLoading(button, true);
            
            const formData = {
                login: document.getElementById('login').value,
                hash: document.getElementById('password').value
            };
            
            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify(formData)
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    showMessage('Login successful! Redirecting...', 'success');
                    setTimeout(() => {
                        window.location.href = data.redirect || '/market';
                    }, 1000);
                } else {
                    showMessage(data.error || 'Login failed', 'error');
                }
            } catch (error) {
                showMessage('Login failed: ' + error.message, 'error');
            } finally {
                setLoading(button, false);
            }
        });

        // Показываем сообщение после регистрации
        const urlParams = new URLSearchParams(window.location.search);
        if (urlParams.get('registered')) {
            showMessage('Registration successful! Please login.', 'success');
        }
    </script>
</body>
</html>
    `, commonStyles, commonScript)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func MarketProductPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Product Details - Market</title>
    %s
    <style>
        .product-container {
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        .product-image {
            width: 100%%;
            max-width: 400px;
            height: auto;
            border-radius: 10px;
            margin-bottom: 20px;
        }
        .product-title {
            font-size: 2rem;
            margin-bottom: 10px;
            color: #333;
        }
        .product-price {
            font-size: 1.5rem;
            color: #2c5aa0;
            font-weight: bold;
            margin-bottom: 15px;
        }
        .product-description {
            font-size: 1.1rem;
            line-height: 1.6;
            color: #666;
            margin-bottom: 20px;
        }
        .product-stock {
            font-size: 1rem;
            color: #28a745;
            margin-bottom: 20px;
        }
        .product-meta {
            font-size: 0.9rem;
            color: #999;
            border-top: 1px solid #eee;
            padding-top: 15px;
        }
        .action-buttons {
            margin-top: 30px;
            display: flex;
            gap: 15px;
        }
        .btn-buy {
            background: #28a745;
            color: white;
            padding: 12px 30px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 1.1rem;
        }
        .btn-buy:hover {
            background: #218838;
        }
        .btn-back {
            background: #6c757d;
            color: white;
            padding: 12px 25px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            text-decoration: none;
            display: inline-block;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="product-container">
            <a href="/market" class="btn-back">← Back to Market</a>
            
            <div id="message"></div>
            
            <div id="product-content">
                <!-- Контент будет загружен через JavaScript -->
                <div style="text-align: center; padding: 50px;">
                    Loading product information...
                </div>
            </div>
        </div>
    </div>

    %s

    <script>
        async function loadProduct() {
            try {
                const productId = '%s';
                const response = await fetch('/api/market/' + productId);
                
                if (!response.ok) {
                    throw new Error('Product not found');
                }
                
                const product = await response.json();
                displayProduct(product);
                
            } catch (error) {
                document.getElementById('product-content').innerHTML = 
                    '<div style="text-align: center; padding: 50px; color: #dc3545;">' +
                    '<h2>Product Not Found</h2>' +
                    '<p>' + error.message + '</p>' +
                    '<a href="/market" class="btn-back">Return to Market</a>' +
                    '</div>';
            }
        }

        function displayProduct(product) {
            const stockStatus = product.amount > 0 ? 
                '✓ In stock: ' + product.amount + ' items' : 
                '✗ Out of stock';
                
            const buyButtonText = product.amount > 0 ? 'Add to Cart' : 'Out of Stock';
            const buyButtonDisabled = product.amount === 0 ? 'disabled' : '';
            
            const content = 
                '<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 40px; align-items: start;">' +
                    '<div>' +
                        '<img src="https://faktodrom.com/i/0000ey00rGGJ/804657_original.png" alt="' + product.name + '" ' +
                        'class="product-image" onerror="this.src=\'/static/placeholder.jpg\'">' +
                    '</div>' +
                    '<div>' +
                        '<h1 class="product-title">' + product.name + '</h1>' +
                        '<div class="product-price">' + formatPrice(product.price) + ' ₽</div>' +
                        '<div class="product-stock">' + stockStatus + '</div>' +
                        '<div class="product-description">' + product.description + '</div>' +
                        '<div class="action-buttons">' +
                            '<button class="btn-buy" ' + buyButtonDisabled + '>' + buyButtonText + '</button>' +
                            '<button class="btn-back" onclick="addToWishlist(' + product.id + ')">' +
                                '♥ Add to Wishlist' +
                            '</button>' +
                        '</div>' +
                        '<div class="product-meta">' +
                            '<strong>Product ID:</strong> ' + product.id + '<br>' +
                            '<strong>Added:</strong> ' + formatDate(product.created_At) +
                        '</div>' +
                    '</div>' +
                '</div>';
            
            document.getElementById('product-content').innerHTML = content;
        }

        function formatPrice(price) {
            return new Intl.NumberFormat('ru-RU').format(price);
        }

        function formatDate(dateString) {
            return new Date(dateString).toLocaleDateString('ru-RU', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            });
        }

        function addToWishlist(productId) {
            showMessage('Added to wishlist!', 'success');
        }

        // Загружаем товар при загрузке страницы
        document.addEventListener('DOMContentLoaded', loadProduct);
    </script>
</body>
</html>
    `, commonStyles, commonScript, id)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
