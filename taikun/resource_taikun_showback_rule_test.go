package taikun

import (
	"fmt"
	"github.com/itera-io/taikungoclient/client/showback"
	"math"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/models"
)

func init() {
	resource.AddTestSweepers("taikun_showback_rule", &resource.Sweeper{
		Name: "taikun_showback_rule",
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := showback.NewShowbackRulesListParams().WithV(ApiVersion)

			var showbackRulesList []*models.ShowbackRulesListDto

			for {
				response, err := apiClient.client.Showback.ShowbackRulesList(params, apiClient)
				if err != nil {
					return err
				}
				showbackRulesList = append(showbackRulesList, response.GetPayload().Data...)
				if len(showbackRulesList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(showbackRulesList))
				params = params.WithOffset(&offset)
			}

			for _, e := range showbackRulesList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					params := showback.NewShowbackDeleteRuleParams().WithV(ApiVersion).WithBody(&models.DeleteShowbackRuleCommand{ID: e.ID})
					_, err = apiClient.client.Showback.ShowbackDeleteRule(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunShowbackRuleConfig = `
resource "taikun_showback_rule" "foo" {
  name = "%s"
  price = %f
  metric_name = "%s"
  type = "%s"
  kind = "%s"
  label {
    key = "key"
    value = "value"
  }
  project_alert_limit = %d
  global_alert_limit = %d
}
`

func TestAccResourceTaikunShowbackRule(t *testing.T) {
	name := randomTestName()
	price := math.Round(rand.Float64()*10000) / 100
	metricName := randomString()
	typeS := []string{"Count", "Sum"}[rand.Int()%2]
	kind := []string{"General", "External"}[rand.Int()%2]
	projectLimit := rand.Int31()
	globalLimit := rand.Int31()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunShowbackRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackRuleConfig,
					name,
					price,
					metricName,
					typeS,
					kind,
					projectLimit,
					globalLimit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunShowbackRuleExists,
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "metric_name", metricName),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "price", fmt.Sprint(price)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "type", typeS),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "kind", kind),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "project_alert_limit", fmt.Sprint(projectLimit)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "global_alert_limit", fmt.Sprint(globalLimit)),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_id"),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_name"),
				),
			},
			{
				ResourceName:      "taikun_showback_rule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunShowbackRuleUpdate(t *testing.T) {
	name := randomTestName()
	price := math.Round(rand.Float64()*10000) / 100
	metricName := randomString()
	typeS := "Count"
	kind := "General"
	projectLimit := rand.Int31()
	globalLimit := rand.Int31()

	newName := randomTestName()
	newPrice := math.Round(rand.Float64()*10000) / 100
	newMetricName := randomString()
	newTypeS := "Sum"
	newKind := "External"
	newProjectLimit := rand.Int31()
	newGlobalLimit := rand.Int31()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunShowbackRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackRuleConfig,
					name,
					price,
					metricName,
					typeS,
					kind,
					projectLimit,
					globalLimit,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunShowbackRuleExists,
					resource.TestCheckResourceAttrSet("taikun_showback_rule.foo", "id"),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "metric_name", metricName),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "price", fmt.Sprint(price)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "type", typeS),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "kind", kind),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "project_alert_limit", fmt.Sprint(projectLimit)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "global_alert_limit", fmt.Sprint(globalLimit)),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_id"),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_name"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackRuleConfig,
					newName,
					newPrice,
					newMetricName,
					newTypeS,
					newKind,
					newProjectLimit,
					newGlobalLimit,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunShowbackRuleExists,
					resource.TestCheckResourceAttrSet("taikun_showback_rule.foo", "id"),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "name", newName),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "metric_name", newMetricName),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "price", fmt.Sprint(newPrice)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "type", newTypeS),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "kind", newKind),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "project_alert_limit", fmt.Sprint(newProjectLimit)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "global_alert_limit", fmt.Sprint(newGlobalLimit)),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_id"),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_name"),
				),
			},
		},
	})
}

const testAccResourceTaikunShowbackRuleWithCredentialsConfig = `
resource "taikun_showback_credential" "foo" {
  name            = "%s"

  password = "%s"
  url = "%s"
  username = "%s"
}

resource "taikun_showback_rule" "foo" {
  name = "%s"
  price = %f
  metric_name = "%s"
  type = "%s"
  kind = "%s"
  label {
    key = "key"
    value = "value"
  }
  project_alert_limit = %d
  global_alert_limit = %d
  showback_credential_id = resource.taikun_showback_credential.foo.id
}
`

func TestAccResourceTaikunShowbackRuleWithCredentials(t *testing.T) {
	showbackCredentialName := randomTestName()
	name := randomTestName()
	price := math.Round(rand.Float64()*10000) / 100
	metricName := randomString()
	typeS := []string{"Count", "Sum"}[rand.Int()%2]
	kind := "External"
	projectLimit := rand.Int31()
	globalLimit := rand.Int31()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunShowbackRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackRuleWithCredentialsConfig,
					showbackCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					name,
					price,
					metricName,
					typeS,
					kind,
					projectLimit,
					globalLimit,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunShowbackRuleExists,
					resource.TestCheckResourceAttrSet("taikun_showback_rule.foo", "id"),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "metric_name", metricName),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "price", fmt.Sprint(price)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "type", typeS),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "kind", kind),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "project_alert_limit", fmt.Sprint(projectLimit)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "global_alert_limit", fmt.Sprint(globalLimit)),
					resource.TestCheckResourceAttrSet("taikun_showback_rule.foo", "showback_credential_id"),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "showback_credential_name", showbackCredentialName),
				),
			},
			{
				ResourceName:      "taikun_showback_rule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTaikunShowbackRuleExists(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_showback_rule" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := showback.NewShowbackRulesListParams().WithV(ApiVersion).WithID(&id)

		response, err := apiClient.client.Showback.ShowbackRulesList(params, apiClient)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("showback rule doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunShowbackRuleDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_showback_rule" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := showback.NewShowbackRulesListParams().WithV(ApiVersion).WithID(&id)

		response, err := apiClient.client.Showback.ShowbackRulesList(params, apiClient)
		if err == nil && response.Payload.TotalCount != 0 {
			return fmt.Errorf("showback rule still exists (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}