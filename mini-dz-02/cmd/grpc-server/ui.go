package main

import (
	"fmt"
	"net/http"
)

// ServeWelcomePage serves the welcome page
func ServeWelcomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, welcomePageHTML)
}

// welcomePageHTML is the HTML template for the welcome page
const welcomePageHTML = `
<html>
	<head>
		<title>Zoo Management System</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
			h1 { color: #333; }
			ul { list-style-type: none; padding: 0; }
			li { margin-bottom: 10px; }
			a { color: #0066cc; text-decoration: none; }
			a:hover { text-decoration: underline; }
		</style>
	</head>
	<body>
		<h1>Zoo Management System API</h1>
		<p>Welcome to the Zoo Management System API. Below are the available endpoints:</p>

		<h2>Animals</h2>
		<ul>
			<li>GET /api/animals - Get all animals</li>
			<li>GET /api/animals/{id} - Get a specific animal</li>
			<li>POST /api/animals - Create a new animal</li>
			<li>DELETE /api/animals/{id} - Delete an animal</li>
			<li>POST /api/animals/{id}/transfer - Transfer an animal to another enclosure</li>
			<li>POST /api/animals/{id}/treat - Treat a sick animal</li>
		</ul>

		<h2>Enclosures</h2>
		<ul>
			<li>GET /api/enclosures - Get all enclosures</li>
			<li>GET /api/enclosures/{id} - Get a specific enclosure</li>
			<li>POST /api/enclosures - Create a new enclosure</li>
			<li>DELETE /api/enclosures/{id} - Delete an enclosure</li>
			<li>GET /api/enclosures/{id}/animals - Get animals in an enclosure</li>
			<li>POST /api/enclosures/{id}/clean - Clean an enclosure</li>
		</ul>

		<h2>Feeding Schedules</h2>
		<ul>
			<li>GET /api/feeding-schedules - Get all feeding schedules</li>
			<li>GET /api/feeding-schedules/{id} - Get a specific feeding schedule</li>
			<li>POST /api/feeding-schedules - Create a new feeding schedule</li>
			<li>DELETE /api/feeding-schedules/{id} - Delete a feeding schedule</li>
			<li>PUT /api/feeding-schedules/{id} - Update a feeding schedule</li>
			<li>POST /api/feeding-schedules/{id}/complete - Mark a feeding schedule as completed</li>
			<li>GET /api/feeding-schedules/due - Get feeding schedules that are due</li>
			<li>GET /api/animals/{id}/feeding-schedules - Get feeding schedules for a specific animal</li>
		</ul>

		<h2>Statistics</h2>
		<ul>
			<li>GET /api/statistics - Get overall zoo statistics</li>
			<li>GET /api/statistics/species - Get animal count by species</li>
			<li>GET /api/statistics/enclosure-utilization - Get enclosure utilization</li>
			<li>GET /api/statistics/health - Get health status statistics</li>
		</ul>
	</body>
</html>
`