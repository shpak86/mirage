package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"mirage/internal/client"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newHttpCommand() *cobra.Command {
	var httpCmd = cobra.Command{
		Use:  "http",
		Long: "Do a single HTTP(s) request",
		Example: `mirage http https://example.com --method get --fp firefox-linux --header "X-Custom-Header:value1" --header "X-Custom-Header:value2" --cookie "session=313373"
cat payload | mirage http https://example.com --method POST --fp chrome-windows --format json`,
		Args: cobra.ExactArgs(1),
		RunE: executeHttp,
	}

	httpCmd.PersistentFlags().StringP("fp", "f", "chrome-android", "fingerprint profile in format PLATFORM-OS. Platforms: chrome, firefox. OS: linux, windows, mac, android, macos. Controls both HTTP headers (UA, client hints) and low-level TLS/JA* fingerprint parameters used for browser impersonation.")
	httpCmd.PersistentFlags().StringP("method", "m", "GET", "HTTP method")
	httpCmd.PersistentFlags().StringSliceP("header", "H", []string{}, "set header in format KEY:VALUE")
	httpCmd.PersistentFlags().StringSliceP("cookie", "C", []string{}, "set cookie")
	httpCmd.PersistentFlags().StringP("output", "o", "resp", "output modes: meta (status + timings), resp (response body), full (status + request + headers + body)")
	httpCmd.PersistentFlags().BoolP("body", "b", false, "read request body from stdio")
	httpCmd.PersistentFlags().StringP("format", "F", "plain", `format output: plain, json`)

	return &httpCmd
}

func executeHttp(cmd *cobra.Command, args []string) error {

	methodFlag, err := cmd.Flags().GetString("method")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return err
	}

	fpFlag, err := cmd.Flags().GetString("fp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return err
	}

	headerFlag, err := cmd.Flags().GetStringSlice("header")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return err
	}
	headers := map[string][]string{}
	for _, it := range headerFlag {
		name, value, found := strings.Cut(it, ":")
		if !found {
			fmt.Fprintf(os.Stderr, "wrong header %s\n", it)
			return fmt.Errorf("wrong header %s", it)
		}
		headerValues, exists := headers[name]
		if !exists {
			headers[name] = []string{value}
		} else {
			headers[name] = append(headerValues, value)
		}
	}

	cookies, err := cmd.Flags().GetStringSlice("cookie")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return err
	}

	var body []byte
	if bodyFlag, err := cmd.Flags().GetBool("body"); err == nil && bodyFlag {
		body, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			return err
		}
	}

	httpClient := client.NewHttpClient()
	request := client.Request{
		Method:      methodFlag,
		Url:         args[0],
		Fingerprint: fpFlag,
		Headers:     headers,
		Cookies:     cookies,
		Body:        body,
	}
	response, err := httpClient.Do(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return err
	}

	format, _ := cmd.Flags().GetString("format")
	printHttpResponse := printHttpResponsePlain
	switch format {
	case `json`:
		printHttpResponse = printHttpResponseJson
	case "plain":
		printHttpResponse = printHttpResponsePlain
	default:
		fmt.Fprintf(os.Stderr, "wrong format %s\n", format)
	}

	output, _ := cmd.Flags().GetString("output")
	switch output {
	case "meta":
		printHttpResponse(response, false, true, false, false)
	case "resp":
		printHttpResponse(response, false, false, false, true)
	case "full":
		printHttpResponse(response, true, true, true, true)
	default:
		fmt.Fprintf(os.Stderr, "error: wrong output mode \"%s\"\n", output)
	}

	return nil
}

func printHttpResponsePlain(resp *client.Response, includeRequest, includeStatus, includeHeaders, includeBody bool) error {
	if includeRequest {
		req := resp.RawRequest
		fmt.Fprintln(os.Stdout, "=== Request ===")
		fmt.Fprintf(os.Stdout, "%s %s\n", req.Method, req.URL)

		if len(req.Header) > 0 {
			fmt.Fprintln(os.Stdout, "\nHeaders:")
			for k, vals := range req.Header {
				for _, v := range vals {
					fmt.Fprintf(os.Stdout, "%s: %s\n", k, v)
				}
			}
		}

		if req.Body != nil {
			if body, err := io.ReadAll(req.Body); err == nil {
				fmt.Fprintln(os.Stdout, "\nBody:")
				fmt.Fprintln(os.Stdout, string(body))
			}
		}
	}

	fmt.Fprintln(os.Stdout, "\n=== Response ===")
	if includeStatus {
		fmt.Fprintf(os.Stdout, "HTTP/%d.%d %d %s\n", resp.RawResponse.ProtoMajor, resp.RawResponse.ProtoMinor, resp.RawResponse.StatusCode, resp.RawResponse.Status)
	}

	if includeHeaders {
		for k, vals := range resp.RawResponse.Header {
			for _, v := range vals {
				fmt.Fprintf(os.Stdout, "%s: %s\n", k, v)
			}
		}
		fmt.Fprintln(os.Stdout)
	}

	if includeBody {
		if body, err := io.ReadAll(resp.RawResponse.Body); err == nil {
			fmt.Fprintln(os.Stdout, string(body))
		}
	}

	return nil
}

func printHttpResponseJson(resp *client.Response, includeRequest, includeStatus, includeHeaders, includeBody bool) error {
	responseOutput := make(map[string]interface{})

	if includeRequest {
		req := resp.RawRequest
		requestOutput := map[string]interface{}{
			"method": req.Method,
			"url":    req.URL.String(),
		}
		if len(req.Header) > 0 {
			headers := make(map[string][]string)
			for k, v := range req.Header {
				headers[k] = v
			}
			requestOutput["headers"] = headers
		}
		if req.Body != nil {
			if body, err := io.ReadAll(req.Body); err == nil {
				requestOutput["body"] = string(body)
				// req.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore body
			}
		}
		responseOutput["request"] = requestOutput
	}

	if includeStatus || includeHeaders || includeBody {
		responseOutput["response"] = map[string]interface{}{}
		if output, ok := responseOutput["response"].(map[string]interface{}); ok {
			if includeStatus {
				output["status"] = resp.RawResponse.Status
				output["statusCode"] = resp.RawResponse.StatusCode
				output["proto"] = fmt.Sprintf("HTTP/%d.%d", resp.RawResponse.ProtoMajor, resp.RawResponse.ProtoMinor)
			}
			if includeHeaders {
				headers := make(map[string][]string)
				for k, v := range resp.RawResponse.Header {
					headers[k] = v
				}
				output["headers"] = headers
			}
			if includeBody {
				if body, err := io.ReadAll(resp.RawResponse.Body); err == nil {
					output["body"] = string(body)
				}
			}
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(responseOutput)
}
