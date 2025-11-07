// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcegroupstaggingapi

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_resourcegroupstaggingapi_tag_resources", name="Tag Resources")
func resourceTagResources() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceTagResourcesCreate,
		ReadWithoutTimeout:   resourceTagResourcesRead,
		UpdateWithoutTimeout: resourceTagResourcesUpdate,
		DeleteWithoutTimeout: resourceTagResourcesDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"resource_arn_list": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			names.AttrTags: tftags.TagsSchema(),
		},
	}
}

func resourceTagResourcesCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ResourceGroupsTaggingAPIClient(ctx)

	resourceARNs := flex.ExpandStringValueSet(d.Get("resource_arn_list").(*schema.Set))
	tags := tftags.New(ctx, d.Get(names.AttrTags).(map[string]interface{}))

	input := &resourcegroupstaggingapi.TagResourcesInput{
		ResourceARNList: resourceARNs,
		Tags:            tags.IgnoreAWS().Map(),
	}

	_, err := conn.TagResources(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "tagging resources: %s", err)
	}

	// Use the resource ARN list as the ID (joined with separator)
	id := strings.Join(resourceARNs, ",")
	d.SetId(id)

	return append(diags, resourceTagResourcesRead(ctx, d, meta)...)
}

func resourceTagResourcesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ResourceGroupsTaggingAPIClient(ctx)

	resourceARNs := strings.Split(d.Id(), ",")

	// We need to read tags for all resources to verify they still exist with the tags
	input := &resourcegroupstaggingapi.GetResourcesInput{
		ResourceARNList: resourceARNs,
	}

	output, err := conn.GetResources(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Resource Groups Tagging API Tag Resources (%s): %s", d.Id(), err)
	}

	if output == nil || len(output.ResourceTagMappingList) == 0 {
		d.SetId("")
		return diags
	}

	// Aggregate all tags from all resources
	allTags := make(map[string]string)
	for _, mapping := range output.ResourceTagMappingList {
		for k, v := range keyValueTags(ctx, mapping.Tags).IgnoreAWS().Map() {
			allTags[k] = v
		}
	}

	d.Set("resource_arn_list", resourceARNs)
	if err := d.Set(names.AttrTags, allTags); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	return diags
}

func resourceTagResourcesUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ResourceGroupsTaggingAPIClient(ctx)

	resourceARNs := flex.ExpandStringValueSet(d.Get("resource_arn_list").(*schema.Set))

	if d.HasChange(names.AttrTags) {
		o, n := d.GetChange(names.AttrTags)
		oldTags := tftags.New(ctx, o)
		newTags := tftags.New(ctx, n)

		// Remove old tags
		if len(oldTags) > 0 {
			input := &resourcegroupstaggingapi.UntagResourcesInput{
				ResourceARNList: resourceARNs,
				TagKeys:         oldTags.IgnoreAWS().Keys(),
			}

			_, err := conn.UntagResources(ctx, input)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "untagging resources: %s", err)
			}
		}

		// Add new tags
		if len(newTags) > 0 {
			input := &resourcegroupstaggingapi.TagResourcesInput{
				ResourceARNList: resourceARNs,
				Tags:            newTags.IgnoreAWS().Map(),
			}

			_, err := conn.TagResources(ctx, input)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "tagging resources: %s", err)
			}
		}
	}

	return append(diags, resourceTagResourcesRead(ctx, d, meta)...)
}

func resourceTagResourcesDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ResourceGroupsTaggingAPIClient(ctx)

	resourceARNs := flex.ExpandStringValueSet(d.Get("resource_arn_list").(*schema.Set))
	tags := tftags.New(ctx, d.Get(names.AttrTags).(map[string]interface{}))

	input := &resourcegroupstaggingapi.UntagResourcesInput{
		ResourceARNList: resourceARNs,
		TagKeys:         tags.IgnoreAWS().Keys(),
	}

	_, err := conn.UntagResources(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "untagging resources: %s", err)
	}

	return diags
}
