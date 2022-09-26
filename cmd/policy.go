package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mysocketio/mysocketctl-go/internal/api/models"
	"github.com/mysocketio/mysocketctl-go/internal/http"
	"github.com/spf13/cobra"
)

// policyCmd represents the policy command
var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Manage your global Policies",
}

// policysListCmd represents the policy ls command
var policysListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List your Policies",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := http.NewClient()

		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		policiesPath := "policies"
		if perPage != 0 {
			if page == 0 {
				page = 1
			}
			policiesPath += fmt.Sprintf("?page_size=%d", perPage)
			policiesPath += fmt.Sprintf("&page=%d", page)
		} else {
			if page != 0 {
				policiesPath += fmt.Sprintf("?page_size=%d", 100)
				policiesPath += fmt.Sprintf("&page=%d", page)
			}
		}

		policys := []models.Policy{}
		err = client.Request("GET", policiesPath, &policys, nil)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Error: %v", err))
		}

		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"Name", "Description", "# Sockets"})

		for _, s := range policys {
			var socketIDs string

			for _, p := range s.SocketIDs {
				if socketIDs == "" {
					socketIDs = socketIDs + ", " + p
				}

			}

			t.AppendRow(table.Row{s.Name, s.Description, len(s.SocketIDs)})
		}
		t.SetStyle(table.StyleLight)
		fmt.Printf("%s\n", t.Render())
	},
}

// policyDeleteCmd represents the policy delete command
var policyDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a policy",
	Run: func(cmd *cobra.Command, args []string) {
		if policyName == "" {
			log.Fatalf("error: invalid policy name")
		}

		policy, err := findPolicyByName(policyName)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Error: %v", err))
		}

		client, err := http.NewClient()
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		err = client.Request("DELETE", "policy/"+policy.ID, nil, nil)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Error: %v", err))
		}

		fmt.Println("Policy deleted")
	},
}

// policyShowCmd represents the policy show command
var policyShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a policy",
	Run: func(cmd *cobra.Command, args []string) {
		if policyName == "" {
			log.Fatalf("error: invalid policy name")
		}

		policy, err := findPolicyByName(policyName)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Error: %v", err))
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"Name", "Description", "# Sockets"})
		t.AppendRow(table.Row{policy.Name, policy.Description, len(policy.SocketIDs)})
		t.SetStyle(table.StyleLight)
		fmt.Printf("%s\n", t.Render())

		jsonData, err := json.MarshalIndent(policy.PolicyData, "", "  ")
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		t = table.NewWriter()
		t.AppendHeader(table.Row{"Policy Data"})
		t.AppendRow(table.Row{string(jsonData)})
		t.SetStyle(table.StyleLight)

		fmt.Printf("%s\n", t.Render())

	},
}

func findPolicyByName(name string) (models.Policy, error) {
	client, err := http.NewClient()

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	policiesPath := "policies/find?name=" + name
	policy := models.Policy{}

	err = client.Request("GET", policiesPath, &policy, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error: %v", err))
	}

	return policy, nil
}

func init() {
	rootCmd.AddCommand(policyCmd)
	policyCmd.AddCommand(policysListCmd)
	policyCmd.AddCommand(policyDeleteCmd)
	policyCmd.AddCommand(policyShowCmd)

	policysListCmd.Flags().Int64Var(&perPage, "per_page", 100, "The number of results to return per page.")
	policysListCmd.Flags().Int64Var(&page, "page", 0, "The page of results to return.")

	policyDeleteCmd.Flags().StringVarP(&policyName, "name", "n", "", "Policy Name")
	policyDeleteCmd.MarkFlagRequired("name")

	policyShowCmd.Flags().StringVarP(&policyName, "name", "n", "", "Policy Name")
	policyShowCmd.MarkFlagRequired("name")
}