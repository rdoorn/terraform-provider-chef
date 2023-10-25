package chef

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceChefNodeDelete() *schema.Resource {
	return &schema.Resource{
		Delete: DeleteNode,

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
