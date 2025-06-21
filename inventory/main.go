package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"math/rand"
	"time"
	"strings"
)

type UserAuth struct {
	RemoteAddress string
	BypassCode    string
	AccessLevel   *http.Cookie
}

func main() {
	// Print current working directory for debugging
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
	} else {
		fmt.Println("Current working directory:", dir)
	}

	// Setup static files handler
	setupStaticFiles()

	http.HandleFunc("/", landingPage)
	http.HandleFunc("/admin", adminDashboard)
	http.HandleFunc("/inventory", stockHandler)
	http.HandleFunc("/contact", contactUsHandler)

	startServer()
}

func setupStaticFiles() {
	staticDir := "./static"
	absPath, _ := filepath.Abs(staticDir)
	fmt.Printf("Configuring static files from: %s\n", absPath)

	// Verify static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		fmt.Printf("ERROR: Static directory not found: %s\n", staticDir)
		return
	}

	// Verify hack.html exists
	hackPath := filepath.Join(staticDir, "hack.html")
	if _, err := os.Stat(hackPath); err == nil {
		fmt.Println("Static files verified successfully")
	} else {
		fmt.Printf("WARNING: hack.html not found: %v\n", err)
	}

	// Configure static file server with cache control
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", 
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			fs.ServeHTTP(w, r)
		}),
	))
}


func landingPage(w http.ResponseWriter, r *http.Request) {
	// Get home URL from environment variable with default fallback
	homeURL := os.Getenv("HOME_BUTTON_URL")
	if homeURL == "" {
		homeURL = "http://localhost:8088"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<title>PinewoodStore Inventory</title>
	<style>
		body { font-family: 'Segoe UI', sans-serif; background: #f8f9fa; margin: 0; padding: 0; }
		.header { background: #2c3e50; color: white; padding: 2rem; text-align: center; position: relative; min-height: 120px; }
		.content { max-width: 1000px; margin: 2rem auto; padding: 2rem; background: white; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
		.btn { background: #3498db; color: white; padding: 0.8rem 1.5rem; border: none; border-radius: 6px; cursor: pointer; font-size: 1rem; transition: all 0.3s; }
		.btn:hover { background: #2980b9; transform: translateY(-2px); }
		.footer { text-align: center; padding: 1.5rem; background: #2c3e50; color: white; }
		.card-container { display: flex; gap: 2rem; margin-top: 2rem; }
		.card { flex: 1; padding: 1.5rem; background: white; border-radius: 8px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
		.home-btn { position: absolute; top: 1rem; left: 1rem; background: #e74c3c; padding: 0.6rem 1.2rem; border-radius: 4px; color: white; text-decoration: none; font-size: 0.9rem; transition: all 0.2s; display: flex; align-items: center; gap: 0.5rem; }
		.home-btn:hover { background: #c0392b; transform: translateY(-1px); box-shadow: 0 2px 5px rgba(0,0,0,0.2); }
		.home-icon { width: 16px; height: 16px; }
		@media (max-width: 768px) {
			.card-container { flex-direction: column; }
			.header { padding-top: 3.5rem; }
			.home-btn { top: 0.5rem; left: 50%%; transform: translateX(-50%%); width: max-content; }
		}
	</style>
</head>
<body>
	<div class="header">
		<a href="%s" class="home-btn">
			<svg class="home-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
				<polyline points="9 22 9 12 15 12 15 22"></polyline>
			</svg>
			Home
		</a>
		<h1>PinewoodStore Inventory System</h1>
		<p>Manage your inventory with ease</p>
	</div>
	<div class="content">
		<div class="card-container">
			<div class="card">
				<h2>Admin Dashboard</h2>
				<p>Access full inventory controls</p>
				<button class="btn" onclick="window.location.href='/admin'">Admin Login</button>
			</div>
			<div class="card">
				<h2>View Inventory</h2>
				<p>Browse current stock</p>
				<button class="btn" onclick="window.location.href='/inventory'">View Items</button>
			</div>
		</div>
		<button class="btn" style="margin-top: 2rem;" onclick="window.location.href='/source'">View Source Code</button>
	</div>
	<div class="footer">
		<p>© 2023 PinewoodStore Inventory System</p>
	</div>
</body>
</html>`, homeURL)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}




func adminDashboard(w http.ResponseWriter, r *http.Request) {
	auth := &UserAuth{}
	auth.AccessLevel, _ = r.Cookie("access_tier")
	auth.RemoteAddress = r.Header.Get("X-Forwarded-For")
	auth.BypassCode = r.Header.Get("pinewoodstore-inventory-team-code")

	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
	<title>Inventory Admin | PinewoodStore</title>
	<style>
		body { font-family: 'Segoe UI', sans-serif; margin: 0; padding: 0; background: #f5f7fa; }
		.navbar { background: #2c3e50; overflow: hidden; display: flex; align-items: center; padding: 0 2rem; height: 60px; }
		.navbar a { color: white; text-align: center; padding: 1rem 1.5rem; text-decoration: none; transition: all 0.3s; line-height: 1; }
		.navbar a:hover { background: #3498db; }
		.navbar .logo {
			font-weight: bold;
			margin-right: auto;
			font-size: 1.5rem;
			letter-spacing: 1px;
			line-height: 1;
			color: white; 
			user-select: none;
		}
		.content { padding: 2rem; }
		.inventory-table { width: 100%%; border-collapse: collapse; margin-top: 1rem; }
		.inventory-table th, .inventory-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #ddd; }
		.inventory-table th { background: #3498db; color: white; }
		.inventory-table tr:hover { background: #f1f1f1; }
		.footer { text-align: center; padding: 1.5rem; background: #2c3e50; color: white; margin-top: 2rem; }
		.stats { display: flex; gap: 1rem; margin: 2rem 0; }
		.stat-card { flex: 1; padding: 1.5rem; background: white; border-radius: 8px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); text-align: center; }
	</style>
</head>
<body>
	<div class="navbar">
		<div class="logo">PinewoodStore Inventory</div>
		<a href="/">Home</a>
		<a href="/admin">Dashboard</a>
		<a href="/inventory">Inventory</a>
		<a href="#">Reports</a>
		<a href="/contact">Contact Us</a>
	</div>
	<div class="content">
		%s
	</div>
	<div class="footer">
		<p>© 2023 PinewoodStore Inventory System - Admin Panel</p>
	</div>
</body>
</html>`, checkPrivileges(auth))
}

func stockHandler(w http.ResponseWriter, r *http.Request) {
	auth := &UserAuth{}
	auth.AccessLevel, _ = r.Cookie("access_tier")
	auth.RemoteAddress = r.Header.Get("X-Forwarded-For")

	// Navbar HTML to reuse
	navbar := `
	<div class="navbar">
		<div class="logo">PinewoodStore Inventory</div>
		<a href="/">Home</a>
		<a href="/admin">Dashboard</a>
		<a href="/inventory">Inventory</a>
		<a href="#">Reports</a>
		<a href="/contact">Contact Us</a>
	</div>`

	// Basic CSS for navbar (same as adminDashboard)
	
    css := `
	<style>
		body { font-family: 'Segoe UI', sans-serif; margin: 0; padding: 0; background: #f5f7fa; }
		.navbar { background: #2c3e50; overflow: hidden; display: flex; align-items: center; padding: 0 2rem; height: 60px; }
		.navbar a { color: white; text-align: center; padding: 1rem 1.5rem; text-decoration: none; transition: all 0.3s; line-height: 1; }
		.navbar a:hover { background: #3498db; }
		.navbar .logo {
			font-weight: bold;
			margin-right: auto;
			font-size: 1.5rem;
			letter-spacing: 1px;
			line-height: 1;
			color: white; 
			user-select: none;
		}
		.header { padding: 2rem; text-align: center; }
		.content { padding: 2rem; max-width: 1200px; margin: 2rem auto; background: white; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
		.search-bar { margin-bottom: 2rem; display: flex; gap: 1rem; }
		.search-input { flex: 1; padding: 0.8rem; border: 1px solid #ddd; border-radius: 6px; font-size: 1rem; }
		.btn { background: #3498db; color: white; border: none; padding: 0.8rem 1.2rem; border-radius: 6px; font-size: 1rem; cursor: pointer; transition: background 0.3s; }
		.btn:hover { background: #2c81ba; }
		.inventory-table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
		.inventory-table th, .inventory-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #ddd; }
		.inventory-table th { background: #3498db; color: white; }
		.inventory-table tr:hover { background: #f1f1f1; }
		.footer { text-align: center; padding: 1.5rem; background: #2c3e50; color: white; margin-top: 2rem; }
		.access-result { padding: 1.5rem; border-radius: 8px; margin-top: 2rem; }
	</style>`

	// Access denied HTML snippet
	accessDenied := `
	<div class="content">
		<div class="access-result" style="background:#ffebee; border-left:4px solid #f44336;">
			<h2 style="color:#d32f2f; margin-top:0;">⚠️ Access Denied</h2>
			<p>You don't have administrator privileges.</p>
			<p>Please contact the Inventory Team with your employee ID to request access.</p>
			<div style="margin-top:1rem; padding:0.8rem; background:#fce4ec; border-radius:4px; font-family: monospace;">
				<small>Debug Info:</small>
				<p>Missing or invalid <code>access_tier</code> cookie</p>
			</div>
		</div>
	</div>`

	// Local access required HTML snippet
	localAccessRequired := `
	<div class="content">
		<div class="access-result" style="background:#e3f2fd; border-left:4px solid #2196f3;">
			<h2 style="color:#1565c0; margin-top:0;">🔒 Local Access Required</h2>
			<p>Admin dashboard requires connection from localhost.</p>
			<div style="margin-top:1rem; padding:0.8rem; background:#e1f5fe; border-radius:4px; font-family: monospace;">
				<small>Debug Info:</small>
				<p>Detected IP: <code>%s</code><br>
				Expected: <code>127.0.0.1</code> or <code>localhost</code><br>
				<em style="color:#666;">(System trusts X-Forwarded-For header)</em></p>
			</div>
		</div>
	</div>`

	if auth.AccessLevel == nil || strings.ToLower(auth.AccessLevel.Value) != "administrator" {
		// Show access denied page with navbar and footer
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Inventory Access Denied | PinewoodStore</title>%s</head>
<body>%s%s
	<div class="footer">
		<p>© 2023 PinewoodStore Inventory System</p>
	</div>
</body>
</html>`, css, navbar, accessDenied)
		return
	}

	ip := auth.RemoteAddress
	if ip == "" {
		ip = r.RemoteAddr
	}
	ipParts := strings.Split(ip, ":")
	clientIP := ipParts[0]

	if clientIP != "127.0.0.1" && clientIP != "localhost" {
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Inventory Access Restricted | PinewoodStore</title>%s</head>
<body>%s`+localAccessRequired+`
	<div class="footer">
		<p>© 2023 PinewoodStore Inventory System</p>
	</div>
</body>
</html>`, css, navbar, clientIP)
		return
	}

	// Authorized - show inventory with navbar and footer
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
	<title>Inventory | PinewoodStore</title>%s
</head>
<body>
	%s
	<div class="header">
		<h1>Current Inventory</h1>
	</div>
	<div class="content">
		<div class="search-bar">
			<input type="text" class="search-input" placeholder="Search inventory...">
			<button class="btn">Search</button>
		</div>
		<table class="inventory-table">
			<thead>
				<tr>
					<th>ID</th>
					<th>Item Name</th>
					<th>Category</th>
					<th>Quantity</th>
					<th>Price</th>
				</tr>
			</thead>
			<tbody>
				<tr><td>1001</td><td>Wooden Chair</td><td>Furniture</td><td>24</td><td>$59.99</td></tr>
				<tr><td>1002</td><td>Oak Table</td><td>Furniture</td><td>12</td><td>$199.99</td></tr>
				<tr><td>1003</td><td>Pine Shelf</td><td>Storage</td><td>18</td><td>$39.99</td></tr>
			</tbody>
		</table>
	</div>
	<div class="footer">
		<p>© 2023 PinewoodStore Inventory System</p>
	</div>
</body>
</html>`, css, navbar)
}


func contactUsHandler(w http.ResponseWriter, r *http.Request) {
	// Seed once per request
	rand.Seed(time.Now().UnixNano())

	// Random placeholders
	emails := []string{"support@pinewoodstore.com", "helpdesk@pinewoodstore.com", "inventory@pinewoodstore.com"}
	phones := []string{"+1-800-555-0199", "+1-888-333-4444", "+1-855-777-1212"}
	ids := []string{"INV-9281", "SUP-4532", "CNT-0049"}

	selectedEmail := emails[rand.Intn(len(emails))]
	selectedPhone := phones[rand.Intn(len(phones))]
	supportID := ids[rand.Intn(len(ids))]

	// Contact Page
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
	<title>Contact Us | PinewoodStore Inventory</title>
	<style>
		body { font-family: 'Segoe UI', sans-serif; margin: 0; padding: 0; background: #f5f7fa; }
		.header { background: #2c3e50; color: white; padding: 2rem; text-align: center; }
		.content { max-width: 800px; margin: 2rem auto; padding: 2rem; background: white; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
		.footer { text-align: center; padding: 1.5rem; background: #2c3e50; color: white; margin-top: 2rem; }
	</style>
</head>
<body>
	<div class="header">
		<h1>Contact Us</h1>
		<p>Need help with PinewoodStore Inventory?</p>
	</div>
	<div class="content">
		<h2>Support Details</h2>
		<p><strong>Email:</strong> %s</p>
		<p><strong>Phone:</strong> %s</p>
		<p><strong>Support ID:</strong> %s</p>
		<p>For inquiries, technical issues, or product support, feel free to contact our team. We’re here to help!</p>
	</div>
	<div class="footer">
		<p>© 2025 PinewoodStore Inventory System</p>
	</div>
</body>
</html>`, selectedEmail, selectedPhone, supportID)
}



func checkPrivileges(auth *UserAuth) string {
    // Case 1: No admin cookie present
    if auth.AccessLevel == nil {
        return `<div class="access-result" style="background:#ffebee; padding:1.5rem; border-radius:8px; border-left:4px solid #f44336; margin-top:1rem;">
            <h2 style="color:#d32f2f; margin-top:0;">⚠️ Access Denied</h2>
            <p>You don't have administrator privileges.</p>
            <p>Please contact the Inventory Team with your employee ID to request access.</p>
            <div style="margin-top:1rem; padding:0.8rem; background:#fce4ec; border-radius:4px;">
                <!-- NEW: Complete Hack button -->
			<button onclick="location.href='/static/leakedcode.html'" 
			style="margin-left:1rem; background:#6a1b9a; color:white; border:none; padding:0.8rem 1.5rem; border-radius:6px; cursor:pointer;">
		💥  Source Code Leaked 
			</button>
            </div>
        </div>`
    }

    // Case 2: Invalid cookie value
    if strings.ToLower(auth.AccessLevel.Value) != "administrator" {
        return `<div class="access-result" style="background:#fff3e0; padding:1.5rem; border-radius:8px; border-left:4px solid #ff9800; margin-top:1rem;">
            <h2 style="color:#ff6f00; margin-top:0;">⚠️ Insufficient Privileges</h2>
            <p>Your account doesn't have admin permissions.</p>
            <div style="margin-top:1rem; padding:0.8rem; background:#fff8e1; border-radius:4px;">
                <small>Debug Info:</small>
                <p style="margin:0.2rem 0; font-family:monospace;">
                    Current cookie: <code>access_tier=` + auth.AccessLevel.Value + `</code><br>
                    Required value: <code>administrator</code>
                </p>
            </div>
        </div>`
    }

    // Extract IP (vulnerable to X-Forwarded-For spoofing)
    clientIP := ""
    if ipParts := strings.Split(auth.RemoteAddress, ":"); len(ipParts) > 0 {
        clientIP = ipParts[0]
    }
    
	


    // Case 3: Valid cookie but wrong IP
    if clientIP != "127.0.0.1" && clientIP != "localhost" {
        return `<div class="access-result" style="background:#e3f2fd; padding:1.5rem; border-radius:8px; border-left:4px solid #2196f3; margin-top:1rem;">
            <h2 style="color:#1565c0; margin-top:0;">🔒 Local Access Required</h2>
            <p>Admin dashboard requires connection from localhost.</p>
            <div style="margin-top:1rem; padding:0.8rem; background:#e1f5fe; border-radius:4px;">
                <small>Debug Info:</small>
                <p style="margin:0.2rem 0; font-family:monospace;">
                    Detected IP: <code>` + clientIP + `</code><br>
                    Expected: <code>127.0.0.1</code> or <code>localhost</code><br>
                    <em style="color:#666;">(System trusts X-Forwarded-For header)</em>
                </p>
            </div>
        </div>`
    }

	
	// Case 4: Invalid cookie value
    if auth.BypassCode != "55569877$$$330003434" {
        return `<div class="access-result" style="background:#fff3e0; padding:1.5rem; border-radius:8px; border-left:4px solid #ff9800; margin-top:1rem;">
            <h2 style="color:#ff6f00; margin-top:0;">⚠️ You are not using the correct Team Code</h2>
            <p>Your account doesn't have admin permissions.</p>
            <div style="margin-top:1rem; padding:0.8rem; background:#fff8e1; border-radius:4px;">
                <small>Debug Info:</small>
                
            </div>
        </div>`
    }

    // Case 5: Successful exploit (cookie + spoofed headers)
return `<div class="admin-content" style="margin-top:1rem;">
    <div style="background:#e8f5e9; padding:1rem; border-radius:6px; border-left:4px solid #4caf50;">
        <h2 style="color:#2e7d32; margin-top:0;">✅ Admin Access Granted</h2>
        <div style="display:grid; grid-template-columns:1fr 1fr; gap:1rem; margin-top:1rem;">
            <div>
			    
                <p><strong>Access Tier:</strong> <code style="background:#c8e6c9; padding:0.2rem 0.4rem;">` + auth.AccessLevel.Value + `</code></p>
                <p><strong>IP Address:</strong> <code style="background:#c8e6c9; padding:0.2rem 0.4rem;">` + auth.RemoteAddress + `</code></p>
				<p><strong>Inventory Team Code:</strong> <code style="background:#c8e6c9; padding:0.2rem 0.4rem;">` + auth.BypassCode + `</code></p>

            </div>
            <div style="font-size:0.9em;">
                <p><strong>Security Note:</strong></p>
                <p>This system trusts client-provided headers for authentication</p>
            </div>
        </div>
    </div>

    <!-- Admin Panel Content -->
    <div style="margin-top:2rem;">
        <h3 style="border-bottom:1px solid #ddd; padding-bottom:0.5rem;">Inventory Controls</h3>
        <div class="admin-stats" style="display:grid; grid-template-columns:repeat(3, 1fr); gap:1rem; margin:1rem 0;">
            <div style="background:#f5f5f5; padding:1rem; border-radius:6px;">
                <h4 style="margin-top:0;">Total Items</h4>
                <p style="font-size:1.5em; margin:0.5rem 0; color:#2e7d32;">54</p>
            </div>
            <div style="background:#f5f5f5; padding:1rem; border-radius:6px;">
                <h4 style="margin-top:0;">Low Stock</h4>
                <p style="font-size:1.5em; margin:0.5rem 0; color:#d32f2f;">7</p>
            </div>
            <div style="background:#f5f5f5; padding:1rem; border-radius:6px;">
                <h4 style="margin-top:0;">Categories</h4>
                <p style="font-size:1.5em; margin:0.5rem 0; color:#1976d2;">12</p>
            </div>
        </div>

        <!-- Existing button -->
        <button style="background:#2e7d32; color:white; border:none; padding:0.8rem 1.5rem; border-radius:6px; cursor:pointer;">
            Manage Inventory
        </button>

        <!-- NEW: Complete Hack button -->
        <button onclick="location.href='/static/hack.html'" 
        style="margin-left:1rem; background:#6a1b9a; color:white; border:none; padding:0.8rem 1.5rem; border-radius:6px; cursor:pointer;">
    💥  Complete Hack
		</button>

    </div>
</div>`
}

func startServer() {
	port := 85
	host := fmt.Sprintf("0.0.0.0:%d", port)
	fmt.Printf("PinewoodStore Inventory Management running on: http://%s\n", host)
	http.ListenAndServe(host, nil)
}