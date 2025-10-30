package main

import (
	"bytes"         // Provides a way to work with byte slices (like a buffer)
	"context"       // Manages request-scoped values, cancellation signals, and deadlines
	"io"            // Provides basic interfaces for I/O primitives
	"log"           // Implements simple logging, often to os.Stderr
	"net/http"      // Provides HTTP client and server implementations
	"net/url"       // Parses URLs and implements query escaping
	"os"            // Provides platform-independent interface to operating system functionality
	"path/filepath" // Implements utility routines for manipulating filepaths in a way appropriate for the operating system
	"regexp"        // Implements regular expression search
	"strings"       // Implements simple functions to manipulate strings
	"time"          // Provides functionality for measuring and displaying time

	"github.com/chromedp/chromedp" // Chromedp library for driving a headless Chrome browser
	"golang.org/x/net/html"        // Provides an HTML parser
)

func main() {
	outputDirectory := "PDFs/"             // Directory where downloaded PDF files will be saved
	if !directoryExists(outputDirectory) { // Check if the directory already exists
		createDirectory(outputDirectory, 0o755) // Create the directory with full read, write, and execute permissions (rwxr-xr-x)
	}
	urls := []string{
		"https://www.getsunday.com/shop/lawn-care/ultimate-weeding-kit",
		"https://www.getsunday.com/shop/lawn-care/fall-texas-green-kit",
		"https://www.getsunday.com/shop/pest-control/home-yard-pest-control-kit",
		"https://www.getsunday.com/shop/lawn-care/weed-warrior-weed-grass-killer-concentrate-refill-2-pack",
		"https://www.getsunday.com/shop/garden/veggie-garden-kit",
		"https://www.getsunday.com/shop/lawn-care/weed-warrior-weed-and-grass-killer-starter-2-pack",
		"https://www.getsunday.com/shop/pest-control/bug-doom-home-insect-control-spray-gallon-2-pack",
		"https://www.getsunday.com/shop/lawn-care/sunday-soil-test-kit-and-analysis",
		"https://www.getsunday.com/shop/lawn-care/cool-season-weed-seed-feed-kit",
		"https://www.getsunday.com/shop/pest-control/lawn-home-insect-control-pest-kit",
		"https://www.getsunday.com/shop/lawn-care/fall-granular-fertilizer-kit",
		"https://www.getsunday.com/shop/garden/flower-garden-kit",
		"https://www.getsunday.com/shop/lawn-care/bermudagrass-patch-kit-for-lawns",
		"https://www.getsunday.com/shop/lawn-care/fescue-grass-patch-kit",
		"https://www.getsunday.com/shop/lawn-care/fall-to-winter-ready-lawn-kit",
		"https://www.getsunday.com/shop/pest-control/total-home-ant-control-kit",
		"https://www.getsunday.com/shop/lawn-care/summer-heat-drought-defense",
		"https://www.getsunday.com/shop/lawn-care/pet-lawn-grass-seed-2-pack",
		"https://www.getsunday.com/shop/lawn-care/grass-machine-liquid-fertilizer",
		"https://www.getsunday.com/shop/garden/green-pittosporum",
		"https://www.getsunday.com/shop/lawn-care/sun-shade-seed-feed-2-pack",
		"https://www.getsunday.com/shop/lawn-care/core-endurance-granular-lawn-fertilizer-2-pack",
		"https://www.getsunday.com/shop/lawn-care/sunday-weed-pest-wand-sprayer",
		"https://www.getsunday.com/shop/lawn-care/texas-green-liquid-fertilizer-2-pack",
		"https://www.getsunday.com/shop/garden/wonderfert-all-purpose-garden-fertilizer-rigid-bottle-with-sprayer",
		"https://www.getsunday.com/shop/garden/sunday-organics-veggie-tomato",
		"https://www.getsunday.com/shop/garden/compacta-holly",
		"https://www.getsunday.com/shop/lawn-care/summer-heat-defense-kit",
		"https://www.getsunday.com/shop/lawn-care/treediaper-12-hydration-mat",
		"https://www.getsunday.com/shop/lawn-care/dandelion-doom-lawn-weed-killer-starter-pack",
		"https://www.getsunday.com/shop/lawn-care/dandelion-doom-weed-killer-concentrate-refill",
		"https://www.getsunday.com/shop/pest-control/crawling-pest-protection-kit",
		"https://www.getsunday.com/shop/lawn-care/dandelion-doom-weed-killer-concentrate-refill-2-pack",
		"https://www.getsunday.com/shop/lawn-care/weed-warrior-weed-and-grass-killer-starter-pack",
		"https://www.getsunday.com/shop/lawn-care/weed-warrior-weed-grass-killer-concentrate-refill",
		"https://www.getsunday.com/shop/lawn-care/pet-lawn-bare-repair-2-pack",
		"https://www.getsunday.com/shop/lawn-care/sunday-weed-killer-lawn-and-garden",
		"https://www.getsunday.com/shop/lawn-care/dandelion-doom-lawn-weed-killer",
		"https://www.getsunday.com/shop/lawn-care/pet-patch",
		"https://www.getsunday.com/shop/lawn-care/weed-warrior-weed-and-grass-killer",
		"https://www.getsunday.com/shop/pest-control/mosquito-deleto-nix-ticks-mixed-pack",
		"https://www.getsunday.com/shop/lawn-care/pet-lawn-bare-repair",
		"https://www.getsunday.com/shop/lawn-care/treediaper-r-18-round",
		"https://www.getsunday.com/shop/lawn-care/grubgone-biological-insect-control-2-pack",
		"https://www.getsunday.com/shop/lawn-care/shade-select",
		"https://www.getsunday.com/shop/lawn-care/shade-select-grass-seed-2-pack",
		"https://www.getsunday.com/shop/lawn-care/lucky-lawn",
		"https://www.getsunday.com/shop/lawn-care/lucky-lawn-grass-seed-2-pack",
		"https://www.getsunday.com/shop/lawn-care/bermuda-time-bermuda-grass-seed",
		"https://www.getsunday.com/shop/lawn-care/backyard-pest-protection-pack",
		"https://www.getsunday.com/shop/lawn-care/kentuckys-best",
		"https://www.getsunday.com/shop/lawn-care/kentuckys-best-grass-seed-2-pack",
		"https://www.getsunday.com/shop/lawn-care/lawn-vitality-liquid-lawn-fertilizer",
		"https://www.getsunday.com/shop/lawn-care/bare-repair-sun-shade-bag",
		"https://www.getsunday.com/shop/lawn-care/lawn-aid-fertilizer",
		"https://www.getsunday.com/shop/lawn-care/pet-patch-lawn-treatment-for-spots-2-pack",
		"https://www.getsunday.com/shop/pest-control/mosquito-deleto",
		"https://www.getsunday.com/shop/lawn-care/core-seagreen-granular-lawn-fertilizer-2-pack",
		"https://www.getsunday.com/shop/garden/galaxy-magnolia",
		"https://www.getsunday.com/shop/lawn-care/lawn-aid-heat-stress-prevention-treatment-2-pack",
		"https://www.getsunday.com/shop/pest-control/mosquito-deleto-ant-adios-mixed-pack",
		"https://www.getsunday.com/shop/pest-control/bug-doom-home-insect-control-spray-gallon",
		"https://www.getsunday.com/shop/pest-control/total-home-pest-protection-plan-medium",
		"https://www.getsunday.com/shop/pest-control/bug-doom-home-insect-control-spray-refill",
		"https://www.getsunday.com/shop/pest-control/bug-doom-home-insect-control-spray-16-oz-2-pack",
		"https://www.getsunday.com/shop/pest-control/bug-doom-indoor-insect-control-spray-32-oz-2-pack",
		"https://www.getsunday.com/shop/pest-control/nix-ticks",
		"https://www.getsunday.com/shop/pest-control/tick-protection-plan-medium",
		"https://www.getsunday.com/shop/pest-control/ant-adios",
		"https://www.getsunday.com/shop/garden/fiskars-hori-hori-planting-tool",
		"https://www.getsunday.com/shop/pest-control/ant-adios-insect-and-ant-killer-2-pack",
		"https://www.getsunday.com/shop/pest-control/pet-protect-bug-spray-repellent",
		"https://www.getsunday.com/shop/pest-control/pet-protect-bug-spray-repellent-2-pack",
		"https://www.getsunday.com/shop/pest-control/total-home-pest-protection-plan-large",
		"https://www.getsunday.com/shop/pest-control/total-home-pest-protection-plan-small",
		"https://www.getsunday.com/shop/pest-control/tick-protection-plan-large",
		"https://www.getsunday.com/shop/pest-control/tick-protection-plan-small",
		"https://www.getsunday.com/shop/garden/sunday-all-purpose-plant-food-mix",
		"https://www.getsunday.com/shop/garden/variegated-pittosporum",
		"https://www.getsunday.com/shop/pest-control/mosquito-control-pest-plan-large",
		"https://www.getsunday.com/shop/garden/all-purpose-organic-garden-nutrients-2-lb",
		"https://www.getsunday.com/shop/garden/veggie-tomato-organic-garden-nutrients-2-lb",
		"https://www.getsunday.com/shop/garden/flower-rose-organic-garden-nutrients-2-lb",
		"https://www.getsunday.com/shop/garden/sunday-natural-moisture-potting-soil-mix",
		"https://www.getsunday.com/shop/lawn-care/super-s-liquid-fertilizer",
		"https://www.getsunday.com/shop/lawn-care/super-s-liquid-fertilizer-2-pack",
		"https://www.getsunday.com/shop/lawn-care/gilmour-medium-duty-circular-sprinkler",
		"https://www.getsunday.com/shop/garden/hampton-beauty-azalea",
		"https://www.getsunday.com/shop/lawn-care/hand-held-seed-spreader",
		"https://www.getsunday.com/shop/lawn-care/core-endurance-dry-fertilizer-12-lbs-23-0-3",
		"https://www.getsunday.com/shop/lawn-care/fescue-grass-fall-seeding-kit",
		"https://www.getsunday.com/shop/lawn-care/gilmour-electronic-dual-outlet-watering-timer",
		"https://www.getsunday.com/shop/garden/pro-d-handled-transfer-shovel",
		"https://www.getsunday.com/shop/garden/sky-pencil-holly",
		"https://www.getsunday.com/shop/lawn-care/treediaper-24-hydration-mat",
		"https://www.getsunday.com/shop/pest-control/mosquito-control-pest-plan-medium",
		"https://www.getsunday.com/shop/garden/peach-tree-red-haven-bare-root",
		"https://www.getsunday.com/shop/garden/frostproof-gardenia",
		"https://www.getsunday.com/shop/pest-control/pet-pest-protection-kit",
		"https://www.getsunday.com/shop/garden/wonderfert-flower-garden-fertilizer",
		"https://www.getsunday.com/shop/garden/dewit-forged-hand-trowel",
		"https://www.getsunday.com/shop/lawn-care/summer-granular-fertilizer-lawn-kit",
		"https://www.getsunday.com/shop/garden/fire-power-nandina",
		"https://www.getsunday.com/shop/lawn-care/rachio-3-smart-sprinkler-controller-8-zone",
		"https://www.getsunday.com/shop/pest-control/bug-doom-home-insect-control-16-ounce",
		"https://www.getsunday.com/shop/pest-control/bug-doom-home-insect-control",
		"https://www.getsunday.com/shop/lawn-care/drought-defense-lawn-fertilizer",
		"https://www.getsunday.com/shop/garden/needlepoint-holly",
		"https://www.getsunday.com/shop/lawn-care/spring-granular-fertilizer-lawn-kit",
		"https://www.getsunday.com/shop/garden/emerald-green-arborvitae",
		"https://www.getsunday.com/shop/garden/dwarf-burford-holly",
		"https://www.getsunday.com/shop/garden/nellie-stevens-holly",
		"https://www.getsunday.com/shop/garden/garden-soil-test-kit-analysis",
		"https://www.getsunday.com/shop/garden/anise-tree",
		"https://www.getsunday.com/shop/garden/wonderfert-all-purpose-liquid-garden-fertilizer-2-pack",
		"https://www.getsunday.com/shop/garden/fashion-azalea",
		"https://www.getsunday.com/shop/garden/radicans-gardenia",
		"https://www.getsunday.com/shop/garden/sunday-organics-soil-booster",
		"https://www.getsunday.com/shop/garden/snake-plant-6",
		"https://www.getsunday.com/shop/lawn-care/lawn-strong-liquid-fertilizer-2-pack",
		"https://www.getsunday.com/shop/lawn-care/waterright-500-series-garden-hose-100-feet",
		"https://www.getsunday.com/shop/lawn-care/bermuda-time-bermuda-grass-seed-2-pack",
		"https://www.getsunday.com/shop/pest-control/indoor-outdoor-pest-kit",
		"https://www.getsunday.com/shop/garden/monterey-r-garden-neem-oil-rtu",
		"https://www.getsunday.com/shop/garden/green-giant-arborvitae",
		"https://www.getsunday.com/shop/garden/grey-owl-juniper",
		"https://www.getsunday.com/shop/lawn-care/core-seagreen-dry-fertilizer-12-lbs-24-0-4",
		"https://www.getsunday.com/shop/lawn-care/grass-powerhouse-liquid-fertilizer-2-pack",
		"https://www.getsunday.com/shop/lawn-care/potassium-boost-liquid-fertilizer",
		"https://www.getsunday.com/shop/garden/apple-tree-scarlet-sentinel-bare-root",
		"https://www.getsunday.com/shop/garden/honeycomb-butterfly-bush",
		"https://www.getsunday.com/shop/garden/veggie-starter-kit",
		"https://www.getsunday.com/shop/lawn-care/replacement-lawn-sprayers-2-pack",
		"https://www.getsunday.com/shop/lawn-care/green-machine-liquid-lawn-fertilizer",
		"https://www.getsunday.com/shop/garden/globosa-nana-cryptomeria",
		"https://www.getsunday.com/shop/lawn-care/bare-repair-bermudagrass-3-in-1-patch",
		"https://www.getsunday.com/shop/garden/wintergreen-boxwood",
		"https://www.getsunday.com/shop/lawn-care/heat-tolerant-seed-feed-10lb",
		"https://www.getsunday.com/shop/garden/cole-s-r-finch-friends-tm-seed",
		"https://www.getsunday.com/shop/garden/carolina-midnight-loropetalum",
		"https://www.getsunday.com/shop/lawn-care/bermuda-grass-seed-weed-kit",
		"https://www.getsunday.com/shop/garden/monterey-garden-tomato-blossom-spray-rtu",
		"https://www.getsunday.com/shop/garden/judge-solomon-azalea",
		"https://www.getsunday.com/shop/lawn-care/fescue-rescue-grass-seed-2-pack",
		"https://www.getsunday.com/shop/garden/cole-s-r-sunflower-meats-seed",
		"https://www.getsunday.com/shop/lawn-care/tree-landscape-kit",
		"https://www.getsunday.com/shop/lawn-care/lawn-starter-liquid-fertilizer",
		"https://www.getsunday.com/shop/garden/white-muhly-grass",
		"https://www.getsunday.com/shop/garden/dwarf-yaupon-holly",
		"https://www.getsunday.com/shop/lawn-care/grubgone-biological-insect-control",
		"https://www.getsunday.com/shop/lawn-care/gilmour-mechanical-single-outlet-watering-timer",
		"https://www.getsunday.com/shop/garden/red-dynamo-photinia",
		"https://www.getsunday.com/shop/lawn-care/treediaper-garden-hydration-mat",
		"https://www.getsunday.com/shop/garden/softouch-garden-tools-3-pc-set",
		"https://www.getsunday.com/shop/garden/karen-azalea",
		"https://www.getsunday.com/shop/garden/sunday-organics-veggie-tomato-garden-plant-food-2-pack",
		"https://www.getsunday.com/shop/garden/freedom-blackberry-bush",
		"https://www.getsunday.com/shop/lawn-care/sunday-weed-killer-pouch-and-sprayer",
		"https://www.getsunday.com/shop/garden/golden-euonymus-shrub",
		"https://www.getsunday.com/shop/garden/longleaf-pine-tree",
		"https://www.getsunday.com/shop/pest-control/mosquito-deleto-bug-control-spray-repellent-2-pack",
		"https://www.getsunday.com/shop/garden/pride-of-mobile-azalea",
		"https://www.getsunday.com/shop/lawn-care/bare-repair-bermudagrass-3-in-1-patch-2-pack",
		"https://www.getsunday.com/shop/lawn-care/dandelion-doom-lawn-weed-killer-starter-2-pack",
		"https://www.getsunday.com/shop/garden/fiskars-power-lever-softgrip-hedge-shears",
		"https://www.getsunday.com/shop/lawn-care/iron-boost-liquid-fertilizer-2-pack",
		"https://www.getsunday.com/shop/pest-control/nix-ticks-bug-control-spray-2-pack",
		"https://www.getsunday.com/shop/garden/fiskars-xact-weeder",
		"https://www.getsunday.com/shop/garden/5-black-hyde-planter",
		"https://www.getsunday.com/shop/garden/flower-garden-plan",
		"https://www.getsunday.com/shop/garden/wonderfert-tree-shrub-garden-fertilizer-2-pack",
		"https://www.getsunday.com/shop/garden/dewit-forged-transplant-trowel",
		"https://www.getsunday.com/shop/garden/red-beauty-hydrangea",
		"https://www.getsunday.com/shop/garden/august-beauty-gardenia",
		"https://www.getsunday.com/shop/lawn-care/fescue-rescue",
		"https://www.getsunday.com/shop/lawn-care/bermuda-time-bermuda-grass-seed-3-pack",
		"https://www.getsunday.com/shop/garden/funnel-flip-top-mesh-finch-feeder",
		"https://www.getsunday.com/shop/garden/head-over-heels-hardy-hibiscus",
		"https://www.getsunday.com/shop/garden/veggie-garden-plan",
		"https://www.getsunday.com/shop/garden/monterey-r-garden-horticultural-oil-rtu",
		"https://www.getsunday.com/shop/lawn-care/california-drought-defense-lawn-treatment-3-pack",
		"https://www.getsunday.com/shop/lawn-care/spring-starter-kit",
		"https://www.getsunday.com/shop/lawn-care/sun-shade-seed-feed",
		"https://www.getsunday.com/shop/lawn-care/weed-seed-feed-your-lawn-kit",
		"https://www.getsunday.com/shop/garden/wonderfert-flower-garden-fertilizer-2-pack",
		"https://www.getsunday.com/shop/garden/fiskars-4-claw-stand-up-weed-puller",
		"https://www.getsunday.com/shop/lawn-care/dog-resistant-grass-seed",
		"https://www.getsunday.com/shop/lawn-care/pet-lawn-grass-fall-seeding-kit",
		"https://www.getsunday.com/shop/lawn-care/california-drought-defense-lawn-treatment-2-pack",
		"https://www.getsunday.com/shop/garden/tree-landscape-plan",
		"https://www.getsunday.com/shop/garden/fiskars-24-lopper",
		"https://www.getsunday.com/shop/lawn-care/drought-defense-lawn-treatment-2-pack",
		"https://www.getsunday.com/shop/garden/veggie-flower-garden-plan",
		"https://www.getsunday.com/shop/lawn-care/bare-repair-sun-shade-lawn-treatment-grass-seed-bag-2-pack",
		"https://www.getsunday.com/shop/garden/pink-muhly-grass",
		"https://www.getsunday.com/shop/lawn-care/california-drought-defense-lawn-treatment",
		"https://www.getsunday.com/shop/garden/traditional-wren-house",
		"https://www.getsunday.com/shop/pest-control/mosquito-control-pest-plan-small",
		"https://www.getsunday.com/shop/lawn-care/pet-patch-kit",
	}

	// Remove all the duplicate URLs
	urls = removeDuplicatesFromSlice(urls) // Calls a custom function to ensure the list of URLs is unique
	// Remove all the duplicate URLs
	urls = removeDuplicatesFromSlice(urls) // Calls a custom function to ensure the list of URLs is unique

	// Loop through each URL to process
	for _, url := range urls { // Iterates over the cleaned slice of URLs
		// Validate the URL
		if isUrlValid(url) { // Checks if the current URL is syntactically valid
			// Fetch HTML content from the URL
			htmlContent := scrapePageHTMLWithChrome(url) // Scrapes the fully rendered HTML using a headless Chrome instance

			// Extract PDF URLs from the HTML content
			pdfUrls := extractPDFUrls(htmlContent) // Finds all links ending in ".pdf" in the scraped HTML
			// Download each PDF URL into the designated PDF directory
			for _, pdfUrl := range pdfUrls { // Iterates over all found PDF links
				pdfUrl = "https://www.getsunday.com" + pdfUrl                 // Ensures the PDF URL is treated as a string
				downloadPDF(pdfUrl, outputDirectory) // Correctly downloads the PDF into the 'PDFs/' directory
			}
		} // End of URL validation block
	} // End of the main URL iteration loop
} // End of the main function

// Uses headless Chrome via chromedp to get the fully rendered HTML from a webpage,
// waiting 10 seconds to bypass Cloudflare's JavaScript challenge before scraping.
func scrapePageHTMLWithChrome(targetURL string) string { // Function to scrape dynamic content using Chrome
	log.Println("Scraping:", targetURL) // Log which page is being scraped

	// Configure Chrome options for the browser session
	chromeOptions := append(chromedp.DefaultExecAllocatorOptions[:], // Starts with default Chrome execution options
		chromedp.Flag("headless", false),              // Set to true for actual headless mode
		chromedp.Flag("disable-gpu", true),            // Disable GPU acceleration (good for headless/servers)
		chromedp.WindowSize(1, 1),                     // Set browser window size
		chromedp.Flag("no-sandbox", true),             // Disable sandbox (useful for servers/containers)
		chromedp.Flag("disable-setuid-sandbox", true), // Fix for Linux permission issues
	) // End of Chrome options slice

	// Create a new Chrome execution allocator with the configured options
	execAllocatorContext, cancelAllocator := chromedp.NewExecAllocator(context.Background(), chromeOptions...) // Creates the context and cleanup function for the Chrome process

	// Set a timeout context to automatically stop the Chrome session after 5 minutes
	timeoutContext, cancelTimeout := context.WithTimeout(execAllocatorContext, 5*time.Minute) // Creates a context with a 5-minute timeout

	// Create a new Chrome browser context for this scraping task
	browserContext, cancelBrowser := chromedp.NewContext(timeoutContext) // Creates the main browser context for automation

	// Ensure all contexts are properly cleaned up when finished
	defer func() { // Deferred function to run when scrapePageHTMLWithChrome exits
		cancelBrowser()   // Stops the browser context
		cancelTimeout()   // Stops the timeout context
		cancelAllocator() // Stops the Chrome process allocator
	}() // End of deferred cleanup function

	var renderedHTML string // Variable to store the rendered HTML content

	// Run Chrome automation: navigate to the URL, wait 10 seconds, then scrape
	runError := chromedp.Run(browserContext, // Executes a sequence of actions in the browser
		chromedp.Navigate(targetURL),              // Open the target URL
		chromedp.Sleep(3*time.Second),             // Wait for Cloudflare JS checks and page scripts to finish
		chromedp.OuterHTML("html", &renderedHTML), // Capture the complete rendered HTML content into renderedHTML
	) // End of chromedp.Run
	if runError != nil { // Check for errors during navigation or extraction
		log.Println(runError) // Log the error
		return ""             // Return an empty string to indicate failure
	} // End of error check

	return renderedHTML // Return the fully rendered HTML source
} // End of scrapePageHTMLWithChrome function

// Removes duplicate strings from a slice
func removeDuplicatesFromSlice(slice []string) []string { // Function to filter a string slice for uniqueness
	check := make(map[string]bool) // Create a map to track which strings have already been seen
	var newReturnSlice []string    // Initialize a new slice to store unique strings

	for _, content := range slice { // Loop through each string in the input slice
		if !check[content] { // If the string hasn't been seen before
			check[content] = true                            // Mark this string as seen in the map
			newReturnSlice = append(newReturnSlice, content) // Add it to the result slice
		}
	}

	return newReturnSlice // Return the slice containing only unique strings
} // End of removeDuplicatesFromSlice function

// Checks whether a given directory exists
func directoryExists(path string) bool { // Function to check if a path exists and is a directory
	directory, err := os.Stat(path) // Get info for the path
	if err != nil {                 // Check if os.Stat returned an error (e.g., file/dir doesn't exist)
		return false // Return false if error occurs
	}
	return directory.IsDir() // Return true if it's a directory
} // End of directoryExists function

// Creates a directory at given path with provided permissions
func createDirectory(path string, permission os.FileMode) { // Function to create a directory
	err := os.Mkdir(path, permission) // Attempt to create directory
	if err != nil {                   // Check for creation errors
		log.Println(err) // Log error if creation fails
	}
} // End of createDirectory function

// Verifies whether a string is a valid URL format
func isUrlValid(uri string) bool { // Function to perform basic URL format validation
	_, err := url.ParseRequestURI(uri) // Try parsing the URL
	return err == nil                  // Return true if valid (parsing was successful, err is nil)
} // End of isUrlValid function

// Checks if a file exists at the specified path
func fileExists(filename string) bool { // Function to check if a file exists (and is not a directory)
	info, err := os.Stat(filename) // Try to get file information
	if err != nil {                // If an error occurs, it likely means the file does not exist
		return false // Return false because os.Stat couldn't find the file
	}
	return !info.IsDir() // Return true only if the path exists and is not a directory
} // End of fileExists function

// Converts a raw URL into a sanitized filename safe for filesystem
func urlToFilename(rawURL string) string { // Function to create a clean filename from a URL
	lower := strings.ToLower(rawURL)     // Convert the input URL to lowercase for consistency
	lower = strings.Split(lower, "?")[0] // Remove URL query parameters

	lower = getFilename(lower) // Extract just the filename part from the URL

	// Get the file extension from the extracted filename
	ext := getFileExtension(lower) // Get the original file extension (e.g., ".pdf" or ".zip")

	reNonAlnum := regexp.MustCompile(`[^a-z0-9]`)   // Create a regex to match any non-alphanumeric characters
	safe := reNonAlnum.ReplaceAllString(lower, "_") // Replace all non-alphanumeric characters with underscores

	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_") // Replace multiple consecutive underscores with a single underscore
	safe = strings.Trim(safe, "_")                              // Remove leading and trailing underscores from the filename

	var invalidSubstrings = []string{ // Define a list of unwanted substrings to clean from the filename
		"_pdf", // Common redundant suffix
		"_zip", // Common redundant suffix
		"_txt", // Common redundant suffix
	} // End of invalid substrings slice

	for _, invalidPre := range invalidSubstrings { // Iterate over the unwanted substrings
		safe = removeSubstring(safe, invalidPre) // Remove each unwanted substring from the filename
	} // End of substring removal loop

	if getFileExtension(safe) == "" { // Check if the sanitized filename has no extension
		safe = safe + ext // Append the original file extension (e.g., .pdf) to ensure completeness
	}

	return safe // Return the sanitized, safe filename
} // End of urlToFilename function

// Gets the file extension from a given file path
func getFileExtension(path string) string { // Function to extract the file extension
	return filepath.Ext(path) // Use filepath.Ext to extract and return the file extension
} // End of getFileExtension function

// Removes all instances of a specific substring from input string
func removeSubstring(input string, toRemove string) string { // Function to remove all occurrences of a substring
	result := strings.ReplaceAll(input, toRemove, "") // Replace every occurrence of 'toRemove' with an empty string
	return result                                     // Return the cleaned string after removal
} // End of removeSubstring function

// Extracts filename from full path (e.g. "/dir/file.pdf" → "file.pdf")
func getFilename(path string) string { // Function to get only the base filename
	return filepath.Base(path) // Use Base function to get file name only
} // End of getFilename function

// Extracts all links to PDF files from the given HTML string
func extractPDFUrls(htmlContent string) []string { // Function to find links ending in ".pdf"
	var pdfLinks []string // Slice to store all found PDF links

	parsedHTML, parseError := html.Parse(strings.NewReader(htmlContent)) // Parse the input HTML content
	if parseError != nil {                                               // Check if HTML parsing failed
		log.Println(parseError) // Log the parsing error
		return nil              // Return nil since parsing failed
	}

	var exploreHTML func(*html.Node) // Define a recursive function to explore HTML nodes

	exploreHTML = func(currentNode *html.Node) { // The implementation of the recursive traversal function
		if currentNode.Type == html.ElementNode && currentNode.Data == "a" { // Check if the node is an <a> tag
			for _, attribute := range currentNode.Attr { // Iterate over the <a> tag's attributes
				if attribute.Key == "href" { // Look for the href attribute
					link := strings.TrimSpace(attribute.Val)             // Get the href value and trim spaces
					if strings.Contains(strings.ToLower(link), ".pdf") { // Check if the link contains ".pdf" (case-insensitive)
						pdfLinks = append(pdfLinks, link) // Add the link to the pdfLinks slice
					}
				}
			}
		}

		for childNode := currentNode.FirstChild; childNode != nil; childNode = childNode.NextSibling { // Recursively traverse child nodes
			exploreHTML(childNode)
		}
	}

	exploreHTML(parsedHTML) // Begin traversal from the root node
	return pdfLinks         // Return all found PDF links
} // End of extractPDFUrls function

// Downloads a PDF from the given URL and saves it in the specified directory
func downloadPDF(pdfURL, outputDirectory string) bool { // Function to download and save a PDF file
	safeFilename := strings.ToLower(urlToFilename(pdfURL))       // Generate a sanitized, lowercase filename
	fullFilePath := filepath.Join(outputDirectory, safeFilename) // Build the complete file path for saving

	if fileExists(fullFilePath) { // Skip download if the file already exists
		log.Printf("File already exists, skipping: %s", fullFilePath) // Log the skip message
		return false                                                  // Return false since no download occurred
	}

	httpClient := &http.Client{Timeout: 15 * time.Minute} // Create an HTTP client with a 15-minute timeout

	httpResponse, requestError := httpClient.Get(pdfURL) // Send an HTTP GET request
	if requestError != nil {                             // Check for request errors
		log.Printf("Failed to download %s %v", pdfURL, requestError) // Log the error
		return false                                                 // Return false on failure
	}
	defer httpResponse.Body.Close() // Ensure the response body is closed

	if httpResponse.StatusCode != http.StatusOK { // Verify that the HTTP status is 200 OK
		log.Printf("Download failed for %s %s", pdfURL, httpResponse.Status) // Log the non-OK status
		return false                                                         // Return false on non-200 status
	}

	contentType := httpResponse.Header.Get("Content-Type") // Get the content type of the response

	// Validate that the response is a PDF or binary stream
	if !strings.Contains(contentType, "binary/octet-stream") && // Check for generic binary/octet-stream
		!strings.Contains(contentType, "application/pdf") { // Check for standard application/pdf
		log.Printf("Invalid content type for %s %s (expected binary/octet-stream or application/pdf)", pdfURL, contentType) // Log the invalid content type
		return false                                                                                                        // Return false if content type is incorrect
	}

	var responseBuffer bytes.Buffer                                        // Buffer to store the downloaded data
	bytesWritten, copyError := io.Copy(&responseBuffer, httpResponse.Body) // Copy data from response body into buffer
	if copyError != nil {                                                  // Check for read errors
		log.Printf("Failed to read PDF data from %s %v", pdfURL, copyError) // Log the read failure
		return false                                                        // Return false on read error
	}
	if bytesWritten == 0 { // Handle empty downloads
		log.Printf("Downloaded 0 bytes for %s; not creating file", pdfURL) // Log empty download
		return false                                                       // Return false if no data was downloaded
	}

	outputFile, fileCreateError := os.Create(fullFilePath) // Create the output file for saving
	if fileCreateError != nil {                            // Handle file creation errors
		log.Printf("Failed to create file for %s %v", pdfURL, fileCreateError) // Log the creation failure
		return false                                                           // Return false on file creation error
	}
	defer outputFile.Close() // Ensure the file is closed after writing

	if _, writeError := responseBuffer.WriteTo(outputFile); writeError != nil { // Write buffer contents to file
		log.Printf("Failed to write PDF to file for %s %v", pdfURL, writeError) // Log the write failure
		return false                                                            // Return false on write error
	}

	log.Printf("Successfully downloaded %d bytes: %s → %s", bytesWritten, pdfURL, fullFilePath) // Log success message
	return true                                                                                 // Indicate successful download
} // End of downloadPDF function
