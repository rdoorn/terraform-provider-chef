package chef

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	chefc "github.com/go-chef/chef"
)

func resourceChefNode() *schema.Resource {
	return &schema.Resource{
		Create: CreateNode,
		Update: UpdateNode,
		Read:   ReadNode,
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
			"environment_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "_default",
			},
			"automatic_attributes_json": {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "{}",
				StateFunc: jsonStateFunc,
			},
			"normal_attributes_json": {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "{}",
				StateFunc: jsonStateFunc,
			},
			"default_attributes_json": {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "{}",
				StateFunc: jsonStateFunc,
			},
			"override_attributes_json": {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "{}",
				StateFunc: jsonStateFunc,
			},
			"run_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:      schema.TypeString,
					StateFunc: runListEntryStateFunc,
				},
			},
		},
	}
}

// CreateNode Creates a Chef Node from the resource definition
func CreateNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	node, err := nodeFromResourceData(d)
	if err != nil {
		return err
	}

	_, err = client.Nodes.Post(*node)
	if err != nil {
		return err
	}

	d.SetId(node.Name)
	return ReadNode(d, meta)
}

// UpdateNode Updates a Chef Node to match the resource definition
func UpdateNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	node, err := nodeFromResourceData(d)
	if err != nil {
		return err
	}

	_, err = client.Nodes.Put(*node)
	if err != nil {
		return err
	}

	d.SetId(node.Name)
	return ReadNode(d, meta)
}

// ReadNode Updates the resource object with data retrieved from Chef
func ReadNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	name := d.Id()

	node, err := client.Nodes.Get(name)
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

	d.Set("name", node.Name)
	d.Set("environment_name", node.Environment)

	automaticAttrJSON, err := json.Marshal(node.AutomaticAttributes)
	if err != nil {
		return err
	}
	d.Set("automatic_attributes_json", string(automaticAttrJSON))

	normalAttrJSON, err := json.Marshal(node.NormalAttributes)
	if err != nil {
		return err
	}
	d.Set("normal_attributes_json", string(normalAttrJSON))

	defaultAttrJSON, err := json.Marshal(node.DefaultAttributes)
	if err != nil {
		return err
	}
	d.Set("default_attributes_json", string(defaultAttrJSON))

	overrideAttrJSON, err := json.Marshal(node.OverrideAttributes)
	if err != nil {
		return err
	}
	d.Set("override_attributes_json", string(overrideAttrJSON))

	runListI := make([]interface{}, len(node.RunList))
	for i, v := range node.RunList {
		runListI[i] = v
	}
	d.Set("run_list", runListI)

	return nil
}

// DeleteNode Deletes a Chef Node
func DeleteNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	name := d.Id()
	err := client.Nodes.Delete(name)

	if err == nil {
		d.SetId("")
	}

	return err
}

func nodeFromResourceData(d *schema.ResourceData) (*chefc.Node, error) {

	node := &chefc.Node{
		Name:        d.Get("name").(string),
		Environment: d.Get("environment_name").(string),
		ChefType:    "node",
		JsonClass:   "Chef::Node",
	}

	var err error

	err = json.Unmarshal(
		[]byte(d.Get("automatic_attributes_json").(string)),
		&node.AutomaticAttributes,
	)
	if err != nil {
		return nil, fmt.Errorf("automatic_attributes_json: %s", err)
	}

	err = json.Unmarshal(
		[]byte(d.Get("normal_attributes_json").(string)),
		&node.NormalAttributes,
	)
	if err != nil {
		return nil, fmt.Errorf("normal_attributes_json: %s", err)
	}

	err = json.Unmarshal(
		[]byte(d.Get("default_attributes_json").(string)),
		&node.DefaultAttributes,
	)
	if err != nil {
		return nil, fmt.Errorf("default_attributes_json: %s", err)
	}

	err = json.Unmarshal(
		[]byte(d.Get("override_attributes_json").(string)),
		&node.OverrideAttributes,
	)
	if err != nil {
		return nil, fmt.Errorf("override_attributes_json: %s", err)
	}

	runListI := d.Get("run_list").([]interface{})
	node.RunList = make([]string, len(runListI))
	for i, vI := range runListI {
		node.RunList[i] = vI.(string)
	}

	return node, nil
}
