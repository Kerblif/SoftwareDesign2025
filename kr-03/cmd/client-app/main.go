package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	orderModels "cbd/internal/orders/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Global variables for API configuration
var (
	apiBaseURL = "http://localhost:8080" // Default API Gateway URL
	isWebMode  = false                   // Flag to indicate if running in web mode
)

// isRunningInWebMode checks if the application is running in web mode
func isRunningInWebMode() bool {
	for _, arg := range os.Args {
		if arg == "-webgl" {
			return true
		}
	}
	return false
}

// getAPIBaseURLForWeb determines the API URL when running in web mode
func getAPIBaseURLForWeb() string {
	// In web mode, the API is assumed to be on the same domain as the web app
	// We can't actually get the current domain in this context, so we'll use a placeholder
	// In a real deployment, the API would be accessible at the same domain
	return "/" // This will make requests relative to the current domain
}

// showAPIConfigDialog shows a dialog to configure the API URL in desktop mode
func showAPIConfigDialog(w fyne.Window, onComplete func()) {
	// Create entries for the API configuration
	hostEntry := widget.NewEntry()
	hostEntry.SetText("localhost")
	hostEntry.SetPlaceHolder("IP address or domain name (e.g., localhost, example.com)")

	portEntry := widget.NewEntry()
	portEntry.SetText("8080")
	portEntry.SetPlaceHolder("Port number (e.g., 8080)")

	protocolSelect := widget.NewSelect([]string{"http", "https"}, nil)
	protocolSelect.SetSelected("http")

	var closeForm func()

	// Create a form with the configuration entries
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Protocol", Widget: protocolSelect},
			{Text: "Host/Domain", Widget: hostEntry},
			{Text: "Port", Widget: portEntry},
		},
		OnSubmit: func() {
			// Validate the host/domain
			if hostEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("host/domain cannot be empty"), w)
				return
			}

			// Validate the port if provided
			if portEntry.Text != "" {
				_, err := strconv.Atoi(portEntry.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("invalid port number: %v", err), w)
					return
				}
			}

			// Construct the API URL
			protocol := protocolSelect.Selected
			host := hostEntry.Text
			port := portEntry.Text

			var apiURL string
			if port != "" {
				apiURL = fmt.Sprintf("%s://%s:%s", protocol, host, port)
			} else {
				apiURL = fmt.Sprintf("%s://%s", protocol, host)
			}

			// Validate the constructed URL
			_, err := url.Parse(apiURL)
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid URL: %v", err), w)
				return
			}

			// Set the API URL
			apiBaseURL = strings.TrimSuffix(apiURL, "/")
			onComplete()

			closeForm()
		},
		OnCancel: func() {
			os.Exit(0)
		},
	}

	// Show the dialog with a more descriptive title
	customDialog := dialog.NewCustomWithoutButtons("Configure API Connection", form, w)

	closeForm = func() {
		customDialog.Hide()
	}

	customDialog.Show()
}

func main() {
	// Check if running in web mode
	isWebMode = isRunningInWebMode()

	// Create a new Fyne application
	a := app.New()

	// Create a new window
	w := a.NewWindow("Order Management System")

	// Set up the API URL based on the mode
	if isWebMode {
		// In web mode, use the same domain
		apiBaseURL = getAPIBaseURLForWeb()
		setupMainUI(w)
	} else {
		// In desktop mode, show a dialog to configure the API URL
		showAPIConfigDialog(w, func() {
			setupMainUI(w)
		})
	}

	// Resize the window
	w.Resize(fyne.NewSize(800, 600))

	// Show and run the application
	w.ShowAndRun()
}

// setupMainUI creates and sets up the main UI
func setupMainUI(w fyne.Window) {
	// Create tabs for different sections
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Dashboard", theme.HomeIcon(), createDashboardTab()),
		container.NewTabItemWithIcon("Orders", theme.DocumentIcon(), createOrdersTab(w)),
		container.NewTabItemWithIcon("Payments", theme.MailSendIcon(), createPaymentsTab(w)),
	)

	// Set the window content
	w.SetContent(tabs)
}

// createDashboardTab creates the dashboard tab content
func createDashboardTab() fyne.CanvasObject {
	welcomeLabel := widget.NewLabelWithStyle(
		"Welcome to Order Management System",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	infoLabel := widget.NewLabel("This application allows you to manage orders and payments.")

	// Create a container with all widgets arranged vertically
	return container.NewVBox(
		welcomeLabel,
		widget.NewSeparator(),
		infoLabel,
		widget.NewLabel("Use the tabs above to navigate between different sections."),
	)
}

// createOrdersTab creates the orders tab content
func createOrdersTab(w fyne.Window) fyne.CanvasObject {
	// Data binding for orders
	var orders []orderModels.Order

	// Create a list to display orders
	ordersList := widget.NewList(
		func() int { return len(orders) },
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(5,
				widget.NewLabel("Order ID"),
				widget.NewLabel("Description"),
				widget.NewLabel("Amount"),
				widget.NewLabel("Status"),
				widget.NewLabel("Created At"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(orders) {
				order := orders[id]
				container := obj.(*fyne.Container)

				// Update the labels with order data
				container.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%d", order.ID))
				container.Objects[1].(*widget.Label).SetText(order.Description)
				container.Objects[2].(*widget.Label).SetText(fmt.Sprintf("$%.2f", order.Amount))
				container.Objects[3].(*widget.Label).SetText(string(order.Status))
				// Format the time without using the time package directly
				createdAt := order.CreatedAt.String()
				if len(createdAt) > 19 {
					createdAt = createdAt[:19]
				}
				container.Objects[4].(*widget.Label).SetText(createdAt)
			}
		},
	)

	// User ID entry for fetching orders
	userIDEntry := widget.NewEntry()
	userIDEntry.SetPlaceHolder("User ID")

	// Create refresh button
	refreshButton := widget.NewButtonWithIcon("Refresh Orders", theme.ViewRefreshIcon(), func() {
		if userIDEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("please enter a user ID"), w)
			return
		}

		fetchIDText := userIDEntry.Text

		userID, err := strconv.ParseUint(fetchIDText, 10, 32)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid user ID: %v", err), w)
			return
		}

		// Fetch orders from the API
		fetchedOrders, err := fetchOrders(uint(userID))
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to fetch orders: %v", err), w)
			return
		}

		// Update the orders list
		orders = fetchedOrders
		ordersList.Refresh()

		ordersList.Resize(fyne.NewSize(ordersList.Size().Width, 400))
	})

	// Create new order form
	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetPlaceHolder("Description")

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("Amount")

	userIDOrderEntry := widget.NewEntry()
	userIDOrderEntry.SetPlaceHolder("User ID")

	createOrderButton := widget.NewButtonWithIcon("Create Order", theme.ContentAddIcon(), func() {
		// Validate inputs
		if descriptionEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("description cannot be empty"), w)
			return
		}

		amount, err := strconv.ParseFloat(amountEntry.Text, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid amount: %v", err), w)
			return
		}

		userID, err := strconv.ParseInt(userIDOrderEntry.Text, 10, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid user ID: %v", err), w)
			return
		}

		// Create order via the API
		order, err := postOrder(userID, amount, descriptionEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to create order: %v", err), w)
			return
		}

		// Create a label for order updates
		updateLabel := widget.NewLabel(fmt.Sprintf("Order #%d created successfully. Waiting for updates...", order.ID))

		// Create a progress bar for visual effect
		progress := widget.NewProgressBarInfinite()

		// Create a container for the popup content with animation effect
		content := container.NewVBox(
			updateLabel,
			progress,
		)

		// Create and show the popup
		popupDialog := dialog.NewCustom("Order Update", "Close", content, w)
		popupDialog.Show()

		// Subscribe to order updates via WebSocket
		msgChan, err := subscribeToOrderUpdates(order.ID)
		if err != nil {
			updateLabel.SetText(fmt.Sprintf("Error subscribing to updates: %v", err))
			return
		}

		// Handle incoming messages in a goroutine
		go func() {
			// Wait for the first message
			msg := <-msgChan

			// Update the label with the message info
			if msg != nil {
				updateLabel.SetText(fmt.Sprintf("Amount: $%.2f, Description: %s",
					msg.Amount, msg.Description))
			} else {
				updateLabel.SetText("No update received or connection closed")
			}

			// Add a small delay for the user to see the update
			time.Sleep(2 * time.Second)

			// Close the popup
			popupDialog.Hide()
		}()

		// Clear the form
		descriptionEntry.SetText("")
		amountEntry.SetText("")

		// Refresh the orders list if the user ID matches
		if userIDOrderEntry.Text == userIDEntry.Text {
			refreshButton.OnTapped()
		}
	})

	// Create form layout
	form := container.NewVBox(
		widget.NewLabel("Create New Order"),
		userIDOrderEntry,
		descriptionEntry,
		amountEntry,
		createOrderButton,
	)

	// Create split container
	split := container.NewHSplit(
		container.NewVSplit(
			container.NewVBox(
				widget.NewLabel("Orders List"),
				userIDEntry,
				refreshButton,
			),
			ordersList,
		),
		container.NewVBox(
			form,
		),
	)
	split.SetOffset(0.7)

	return split
}

// createPaymentsTab creates the payments tab content
func createPaymentsTab(w fyne.Window) fyne.CanvasObject {
	// Create balance display
	balanceValue := binding.NewFloat()
	balanceValue.Set(0.0)

	balanceDisplay := widget.NewLabelWithData(binding.FloatToStringWithFormat(balanceValue, "Balance: $%.2f"))

	// Create account ID input
	accountIDEntry := widget.NewEntry()
	accountIDEntry.SetPlaceHolder("Enter User ID")

	// Status indicator
	statusLabel := widget.NewLabel("")

	// Create check balance button
	checkBalanceButton := widget.NewButtonWithIcon("Check Balance", theme.SearchIcon(), func() {
		if accountIDEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("please enter a user ID"), w)
			return
		}

		userID, err := strconv.ParseInt(accountIDEntry.Text, 10, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid user ID: %v", err), w)
			return
		}

		// Check balance via the API
		balance, err := getBalance(userID)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to get balance: %v", err), w)
			statusLabel.SetText("Error: " + err.Error())
			return
		}

		// Update the balance display
		balanceValue.Set(balance)
		statusLabel.SetText("Balance retrieved successfully")
	})

	// Create new account button
	createAccountButton := widget.NewButtonWithIcon("Create Account", theme.ContentAddIcon(), func() {
		if accountIDEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("please enter a user ID"), w)
			return
		}

		userID, err := strconv.ParseInt(accountIDEntry.Text, 10, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid user ID: %v", err), w)
			return
		}

		// Create account via the API
		err = createNewAccount(userID)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to create account: %v", err), w)
			statusLabel.SetText("Error: " + err.Error())
			return
		}

		dialog.ShowInformation("Success", "Account created successfully", w)
		statusLabel.SetText("Account created successfully")

		// Check the balance to confirm
		checkBalanceButton.OnTapped()
	})

	// Create deposit form
	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("Amount")

	depositButton := widget.NewButtonWithIcon("Deposit", theme.ContentAddIcon(), func() {
		if accountIDEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("please enter a user ID"), w)
			return
		}

		amount, err := strconv.ParseFloat(amountEntry.Text, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid amount: %v", err), w)
			return
		}

		userID, err := strconv.ParseInt(accountIDEntry.Text, 10, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid user ID: %v", err), w)
			return
		}

		// Deposit via the API
		err = addDeposit(userID, amount)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to add deposit: %v", err), w)
			statusLabel.SetText("Error: " + err.Error())
			return
		}

		dialog.ShowInformation("Success", fmt.Sprintf("Deposited $%.2f successfully", amount), w)
		statusLabel.SetText(fmt.Sprintf("Deposited $%.2f successfully", amount))

		// Clear the amount entry
		amountEntry.SetText("")

		// Check the balance to confirm
		checkBalanceButton.OnTapped()
	})

	// Create a wider container for the account ID entry to prevent horizontal scrolling
	accountIDContainer := container.NewGridWithColumns(1, accountIDEntry)

	return container.NewVBox(
		widget.NewLabel("Account Management"),
		accountIDContainer, // Place the account ID entry in its own row to give it full width
		container.NewHBox(
			checkBalanceButton,
			createAccountButton,
		),
		balanceDisplay,
		statusLabel,
		widget.NewSeparator(),
		widget.NewLabel("Deposit Funds"),
		amountEntry,
		depositButton,
	)
}
