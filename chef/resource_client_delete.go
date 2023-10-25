package chef

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceChefClientDelete() *schema.Resource {
	return &schema.Resource{
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
		},
	}
}
