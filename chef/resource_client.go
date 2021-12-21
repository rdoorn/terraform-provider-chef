package chef

import (
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	chefc "github.com/go-chef/chef"
)

func resourceChefClient() *schema.Resource {
	return &schema.Resource{
		Create: CreateApiClient,
		Read:   ReadApiClient,
		Delete: DeleteApiClient,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"validator": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
				StateFunc: CryptoHashingStateFunc,
			},
		},
	}
}

func CryptoHashingStateFunc(privateKey interface{}) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(privateKey.(string))))
}

// CreateApiClient Creates a Chef Client from the resource definition
func CreateApiClient(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	apiClient, err := apiClientFromResourceData(d)
	if err != nil {
		return err
	}

	result, err := client.Clients.Create(apiClient.Name, apiClient.Admin)
	if err != nil {
		return err
	}

	d.SetId(apiClient.Name)
	d.Set("private_key", result.PrivateKey)
	return ReadApiClient(d, meta)
}

// ReadApiClient Updates the resource object with data retrieved from Chef
func ReadApiClient(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	name := d.Id()

	apiClient, err := client.Clients.Get(name)
	if err != nil {
		if errRes, ok := err.(*chefc.ErrorResponse); ok {
			if errRes.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		} else {
			return err
		}
	}

	d.Set("name", apiClient.Name)
	d.Set("validator", apiClient.Validator)
	d.Set("admin", apiClient.Admin)
	d.Set("public_key", apiClient.PublicKey)

	return nil
}

// DeleteApiClient Deletes a Chef Client
func DeleteApiClient(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	name := d.Id()
	err := client.Clients.Delete(name)

	if err == nil {
		d.SetId("")
	}

	return err
}

func apiClientFromResourceData(d *schema.ResourceData) (*chefc.ApiClient, error) {

	apiClient := &chefc.ApiClient{
		Name:      d.Get("name").(string),
		PublicKey: d.Get("public_key").(string),
		Admin:     d.Get("admin").(bool),
		ChefType:  "client",
		JsonClass: "Chef::ApiClient",
	}

	return apiClient, nil
}
