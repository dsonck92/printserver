logging {
  development = false
}

http {
  address = "0.0.0.0:8000"
  debug = false
  rate_limit = 20
}

printer {
  lp_binary = "lp"
  lp_args = []
  dest_dir = "/tmp/printer"
}

scanner {
  scanimage_binary = "/usr/bin/scanimage"
  scanimage_args = []
  dest_dir = "/tmp/scanner"
}
