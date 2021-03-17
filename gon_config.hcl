source = ["./build/terraform-provider-chef-darwin-amd64"]
bundle_id = "com.terrycain.github.terraform-provider-chef"

apple_id {
  username = "terry@dolphincorp.co.uk"
  password = "@env:APPLE_APP_PW"
}

sign {
  application_identity = "Developer ID Application: Terry Cain (UT7M7Z36B6)"
}


zip {
  output_path = "build/terraform-provider-chef-darwin-amd64.zip"
}
