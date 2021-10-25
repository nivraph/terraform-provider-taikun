package taikun

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/projects"
	"github.com/itera-io/taikungoclient/models"
)

func init() {
	resource.AddTestSweepers("taikun_project", &resource.Sweeper{
		Name:         "taikun_project",
		Dependencies: []string{},
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := projects.NewProjectsListParams().WithV(ApiVersion)

			var projectList []*models.ProjectListForUIDto

			for {
				response, err := apiClient.client.Projects.ProjectsList(params, apiClient)
				if err != nil {
					return err
				}
				projectList = append(projectList, response.GetPayload().Data...)
				if len(projectList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(projectList))
				params = params.WithOffset(&offset)
			}

			for _, e := range projectList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					params := projects.NewProjectsDeleteParams().WithV(ApiVersion).WithBody(&models.DeleteProjectCommand{ProjectID: e.ID})
					_, _, err = apiClient.client.Projects.ProjectsDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunProjectConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  enable_auto_upgrade = %t
  enable_monitoring = %t
  expiration_date = "%s"
}
`

func TestAccResourceTaikunProject(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	enableAutoUpgrade := true
	enableMonitoring := false
	expirationDate := "01/04/2999"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
				),
			},
			{
				ResourceName:      "taikun_project.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunProjectExtendLifetime(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	enableAutoUpgrade := true
	enableMonitoring := false
	expirationDate := "01/04/2999"
	newExpirationDate := "07/02/3000"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					newExpirationDate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", newExpirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
				),
			},
		},
	})
}

func TestAccResourceTaikunProjectToggleMonitoring(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	enableAutoUpgrade := true
	enableMonitoring := true
	disableMonitoring := false
	expirationDate := "01/04/2999"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					enableAutoUpgrade,
					disableMonitoring,
					expirationDate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_monitoring", fmt.Sprint(disableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "enable_monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
				),
			},
		},
	})
}

const testAccResourceTaikunProjectToggleBackupConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"
}

resource "taikun_backup_credential" "foo" {
  name            = "%s"

  s3_endpoint = "%s"
  s3_region   = "%s"
}

resource "taikun_backup_credential" "foo2" {
  name            = "%s"

  s3_endpoint = "%s"
  s3_region   = "%s"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  %s
}
`

func TestAccResourceTaikunProjectToggleBackup(t *testing.T) {
	cloudCredentialName := randomTestName()
	backupCredentialName := randomTestName()
	backupCredentialName2 := randomTestName()
	projectName := "myproject"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckS3(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectToggleBackupConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					backupCredentialName2,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					projectName,
					"backup_credential_id = resource.taikun_backup_credential.foo.id",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttrSet("taikun_project.foo", "enable_auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "enable_monitoring"),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
					resource.TestCheckResourceAttrPair("taikun_project.foo", "backup_credential_id", "taikun_backup_credential.foo", "id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectToggleBackupConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					backupCredentialName2,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					projectName,
					"backup_credential_id = resource.taikun_backup_credential.foo2.id",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttrSet("taikun_project.foo", "enable_auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "enable_monitoring"),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
					resource.TestCheckResourceAttrPair("taikun_project.foo", "backup_credential_id", "taikun_backup_credential.foo2", "id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectToggleBackupConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					backupCredentialName2,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					projectName,
					"",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttrSet("taikun_project.foo", "enable_auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "enable_monitoring"),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
					resource.TestCheckResourceAttr("taikun_project.foo", "backup_credential_id", ""),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectToggleBackupConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					backupCredentialName2,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					projectName,
					"backup_credential_id = resource.taikun_backup_credential.foo.id",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttrSet("taikun_project.foo", "enable_auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "enable_monitoring"),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
					resource.TestCheckResourceAttrPair("taikun_project.foo", "backup_credential_id", "taikun_backup_credential.foo", "id"),
				),
			},
		},
	})
}

func testAccCheckTaikunProjectExists(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := projects.NewProjectsListParams().WithV(ApiVersion).WithID(&id)

		response, err := apiClient.client.Projects.ProjectsList(params, apiClient)
		if err != nil || len(response.Payload.Data) != 1 {
			return fmt.Errorf("project doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunProjectDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := projects.NewProjectsListParams().WithV(ApiVersion).WithID(&id)

		response, err := apiClient.client.Projects.ProjectsList(params, apiClient)
		if err == nil && len(response.Payload.Data) != 0 {
			return fmt.Errorf("project still exists (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}
