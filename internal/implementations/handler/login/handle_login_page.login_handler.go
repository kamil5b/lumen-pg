package login

import (
	"net/http"
)

func (h *LoginHandlerImplementation) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	// Render login page HTML
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<title>Login - Lumen PG</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
			margin: 0;
			background: #f5f5f5;
		}
		.login-container {
			background: white;
			padding: 40px;
			border-radius: 8px;
			box-shadow: 0 2px 10px rgba(0,0,0,0.1);
			width: 100%;
			max-width: 400px;
		}
		h1 {
			margin-top: 0;
			text-align: center;
		}
		.form-group {
			margin-bottom: 20px;
		}
		label {
			display: block;
			margin-bottom: 5px;
			font-weight: bold;
		}
		input[type="text"],
		input[type="password"] {
			width: 100%;
			padding: 10px;
			border: 1px solid #ddd;
			border-radius: 4px;
			box-sizing: border-box;
		}
		button {
			width: 100%;
			padding: 12px;
			background: #007bff;
			color: white;
			border: none;
			border-radius: 4px;
			cursor: pointer;
			font-size: 16px;
		}
		button:hover {
			background: #0056b3;
		}
		.error {
			color: #dc3545;
			margin-top: 10px;
		}
	</style>
</head>
<body>
	<div class="login-container">
		<h1>Lumen PG</h1>
		<form method="POST" action="/login">
			<div class="form-group">
				<label for="username">Username:</label>
				<input type="text" id="username" name="username" required>
			</div>
			<div class="form-group">
				<label for="password">Password:</label>
				<input type="password" id="password" name="password" required>
			</div>
			<button type="submit">Login</button>
		</form>
	</div>
</body>
</html>`))
}
