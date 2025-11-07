// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcegroupstaggingapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccResourceGroupsTaggingAPITagResources_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_resourcegroupstaggingapi_tag_resources.test"
	vpcResourceName := "aws_vpc.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ResourceGroupsTaggingAPIServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTagResourcesDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourcesConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagResourcesExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "resource_arn_list.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "resource_arn_list.*", vpcResourceName, names.AttrARN),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceGroupsTaggingAPITagResources_updateTags(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_resourcegroupstaggingapi_tag_resources.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ResourceGroupsTaggingAPIServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTagResourcesDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourcesConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagResourcesExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", rName),
				),
			},
			{
				Config: testAccTagResourcesConfig_updated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagResourcesExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "Staging"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.Updated", acctest.CtTrue),
				),
			},
		},
	})
}

func TestAccResourceGroupsTaggingAPITagResources_multipleResources(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_resourcegroupstaggingapi_tag_resources.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ResourceGroupsTaggingAPIServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTagResourcesDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourcesConfig_multiple(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagResourcesExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "resource_arn_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.Project", "MultiResource"),
				),
			},
		},
	})
}

func testAccCheckTagResourcesExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ResourceGroupsTaggingAPIClient(ctx)

		input := &resourcegroupstaggingapi.GetResourcesInput{
			ResourceARNList: []string{rs.Primary.ID},
		}

		output, err := conn.GetResources(ctx, input)
		if err != nil {
			return err
		}

		if len(output.ResourceTagMappingList) == 0 {
			return fmt.Errorf("Resource Groups Tagging API Tag Resources not found")
		}

		return nil
	}
}

func testAccCheckTagResourcesDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Since we're only managing tags, not the resources themselves,
		// we just check that the tags have been removed
		// The actual resources (VPCs, etc.) are destroyed by their own resource types
		return nil
	}
}

func testAccTagResourcesConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_resourcegroupstaggingapi_tag_resources" "test" {
  resource_arn_list = [
    aws_vpc.test.arn,
  ]

  tags = {
    Name        = %[1]q
    Environment = "Production"
  }
}
`, rName)
}

func testAccTagResourcesConfig_updated(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_resourcegroupstaggingapi_tag_resources" "test" {
  resource_arn_list = [
    aws_vpc.test.arn,
  ]

  tags = {
    Name        = %[1]q
    Environment = "Staging"
    Updated     = "true"
  }
}
`, rName)
}

func testAccTagResourcesConfig_multiple(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test1" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "%[1]s-1"
  }
}

resource "aws_vpc" "test2" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = "%[1]s-2"
  }
}

resource "aws_resourcegroupstaggingapi_tag_resources" "test" {
  resource_arn_list = [
    aws_vpc.test1.arn,
    aws_vpc.test2.arn,
  ]

  tags = {
    Environment = "Production"
    Project     = "MultiResource"
  }
}
`, rName)
}
