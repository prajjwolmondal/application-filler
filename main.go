package main

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

func main() {
	log.Println("Starting application")

	pw := setupPlaywright()
	browser := setupFirefox(pw)
	page := createPage(browser)

	url := "https://boards.greenhouse.io/simpplr/jobs/5241142004"
	_, err := page.Goto(url)
	if err != nil {
		log.Fatalf("couldn't go to URL %s | err: %v", url, err)
	}

	fillInLeverJobApplication(page)

	closeFirefoxAndPlaywright(browser, pw)
}

func setupPlaywright() *playwright.Playwright {
	pw, err := playwright.Run()

	if err != nil {
		log.Fatalf("coudln't start playwright: %v", err)
	}

	log.Println("Finished setting up Playwright")
	return pw
}

func setupFirefox(pw *playwright.Playwright) playwright.Browser {
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})

	if err != nil {
		log.Fatalf("couldn't launch Firefox: %v", err)
	}

	log.Println("Finished setting up Firefox")
	return browser

}

func createPage(browser playwright.Browser) playwright.Page {
	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("couldn't create context: %v", err)
	}

	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("couldn't create page: %v", err)
	}

	log.Println("Finished creating page")
	return page
}

func closeFirefoxAndPlaywright(browser playwright.Browser, pw *playwright.Playwright) {
	err := browser.Close()
	if err != nil {
		log.Fatalf("couldn't close Firefox: %v", err)
	}

	err = pw.Stop()
	if err != nil {
		log.Fatalf("couldn't stop Playwright: %v", err)
	}

	log.Println("Closed Firefox and Playwright")
}

func fillFieldIfPresent(locator playwright.Locator, value, fieldNameForLogMessage string) {
	if isPresent, err := locator.IsVisible(); err != nil {
		log.Printf("---Couldn't find %s field, skipping it---", fieldNameForLogMessage)
	} else if isPresent {
		locator.Fill(value)
	}
}

func fillInLeverJobApplication(page playwright.Page) {

	fillFieldIfPresent(page.Locator("#first_name"), "Prajjwol", "first name")
	fillFieldIfPresent(page.Locator("#last_name"), "Mondal", "last name")
	fillFieldIfPresent(page.Locator("#email"), "prajj.mondal@gmail.com", "email")
	fillFieldIfPresent(page.Locator("#phone"), "416-778-9993", "phone")
	fillFieldIfPresent(page.Locator("[autocomplete=\"custom-question-linkedin-profile\"]"), "https://www.linkedin.com/in/prajjwolmondal/", "LinkedIn")
	fillFieldIfPresent(page.Locator("[autocomplete=\"custom-question-website-or-portfolio-link\"]"), "https://prajjwolmondal.github.io/", "Website/portfolio link")

	//TODO: Figure out how to attach resume

	// Following didn't seem to work but I didn't check to see if any error was logged
	// err := page.Locator("div[data-field=\"resume\"] button[data-source=\"attach\"]").SetInputFiles(playwright.InputFile{Name: "Resume [Prajjwol Mondal].pdf"})
	// if err != nil {
	// 	log.Fatal("Unable to attach resume: %v", err)
	// }

	//Following is an attempt at handling the file upload asynchronously but I'm running into
	// fatal error: all goroutines are asleep - deadlock!
	// done := make(chan bool)

	// page.On("filechooser", func(fileChooser playwright.FileChooser) {
	// 	log.Println("File chooser triggered")
	// 	err := fileChooser.SetFiles("Resume [Prajjwol Mondal].pdf")
	// 	if err != nil {
	// 		log.Fatalf("could not set file: %v", err)
	// 	}
	// 	done <- true
	// })

	// err := page.Locator("div[data-field=\"resume\"] button[data-source=\"attach\"]").Click()

	// if err != nil {
	// 	log.Fatalf("could not click attach button: %v", err)
	// }

	// <-done

	log.Println("Successfully filled out application")
}
