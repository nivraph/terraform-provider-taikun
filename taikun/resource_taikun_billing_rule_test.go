package taikun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/prometheus"
	"github.com/itera-io/taikungoclient/models"
	"os"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("taikun_billing_rule", &resource.Sweeper{
		Name: "taikun_billing_rule",
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion)

			var billingRulesList []*models.PrometheusRuleListDto
			for {
				response, err := apiClient.client.Prometheus.PrometheusListOfRules(params, apiClient)
				if err != nil {
					return err
				}
				billingRulesList = append(billingRulesList, response.GetPayload().Data...)
				if len(billingRulesList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(billingRulesList))
				params = params.WithOffset(&offset)
			}

			for _, e := range billingRulesList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					params := prometheus.NewPrometheusDeleteParams().WithV(ApiVersion).WithID(e.ID)
					_, err = apiClient.client.Prometheus.PrometheusDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunBillingRule = `
resource "taikun_billing_credential" "foo" {
  name            = "%s"
  is_locked       = false

  prometheus_password = "%s"
  prometheus_url = "%s"
  prometheus_username = "%s"
}

resource "taikun_billing_rule" "foo" {
  name            = "%s"
  metric_name     =   "coredns_forward_request_duration_seconds"
  price = 1
  type = "Sum"
  billing_credential_id = resource.taikun_billing_credential.foo.id
  label {
    key = "key"
    value = "value"
  }
}
`

func TestAccResourceTaikunBillingRule(t *testing.T) {
	credName := randomTestName()
	ruleName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingRule,
					credName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					ruleName,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunBillingRuleExists,
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "name", ruleName),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "metric_name", "coredns_forward_request_duration_seconds"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "type", "Sum"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "price", "1"),
					resource.TestCheckResourceAttrSet("taikun_billing_rule.foo", "billing_credential_id"),
				),
			},
		},
	})
}

func TestAccResourceTaikunBillingRuleRename(t *testing.T) {
	credName := randomTestName()
	ruleName := randomTestName()
	ruleNameNew := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingRule,
					credName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					ruleName,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunBillingRuleExists,
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "name", ruleName),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "metric_name", "coredns_forward_request_duration_seconds"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "type", "Sum"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "price", "1"),
					resource.TestCheckResourceAttrSet("taikun_billing_rule.foo", "billing_credential_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingRule,
					credName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					ruleNameNew,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunBillingRuleExists,
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "name", ruleNameNew),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "metric_name", "coredns_forward_request_duration_seconds"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "type", "Sum"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "price", "1"),
					resource.TestCheckResourceAttrSet("taikun_billing_rule.foo", "billing_credential_id"),
				),
			},
		},
	})
}

func testAccCheckTaikunBillingRuleExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_rule" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.Prometheus.PrometheusListOfRules(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("access profile doesn't exist")
		}
	}

	return nil
}

func testAccCheckTaikunBillingRuleDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_rule" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.Prometheus.PrometheusListOfRules(params, client)
		if err == nil && response.Payload.TotalCount != 0 {
			return fmt.Errorf("access profile still exists")
		}
	}

	return nil
}