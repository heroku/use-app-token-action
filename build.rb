#!/usr/bin/env ruby
# frozen_string_literal: true

require 'digest'
require 'fileutils'

build_configs = [
  {
    target_os_arch: 'GOOS=darwin GOARCH=amd64',
    dest_os_arch: 'macOS-X64'
  },
  {
    target_os_arch: 'GOOS=darwin GOARCH=arm64',
    dest_os_arch: 'macOS-ARM64'
  },
  {
    target_os_arch: 'GOOS=linux GOARCH=amd64',
    dest_os_arch: 'Linux-X64'
  },
  {
    target_os_arch: 'GOOS=linux GOARCH=arm64',
    dest_os_arch: 'Linux-ARM64'
  },
  {
    target_os_arch: 'GOOS=windows GOARCH=amd64',
    dest_os_arch: 'Windows-X64'
  },
  {
    target_os_arch: 'GOOS=windows GOARCH=arm64',
    dest_os_arch: 'Windows-ARM64'
  }
]

build_configs.each do |cfg|
  tmp_binary_path = "tmp/main-#{cfg[:dest_os_arch]}"
  bin_binary_path = "bin/main-#{cfg[:dest_os_arch]}"
  build_cmd = "#{cfg[:target_os_arch]} go build -ldflags=\"-s -w\" -o #{tmp_binary_path} cmd/main.go"

  puts("Building binary for #{cfg[:dest_os_arch]}")
  build_cmd_result = system(build_cmd)

  if build_cmd_result != true
    raise "An error occurred while creating the binary for #{cfg[:dest_os_arch]}: #{tmp_binary_path}"
  end

  if File.exist?(bin_binary_path)
    tmp_check_sum = Digest::SHA1.file(tmp_binary_path).hexdigest
    bin_check_sum = Digest::SHA1.file(bin_binary_path).hexdigest

    if tmp_check_sum == bin_check_sum
      puts("  Binary for #{cfg[:dest_os_arch]} is up to date. 👍")
    else
      puts("  Replacing binary for #{cfg[:dest_os_arch]} with newer version")
      puts("  SHA256: #{tmp_check_sum} --> ✅  SHA256: #{bin_check_sum}")
      FileUtils.mv(tmp_binary_path, bin_binary_path, force: true)
    end
  else
    puts("  Adding binary for #{cfg[:dest_os_arch]}")
    FileUtils.mv(tmp_binary_path, bin_binary_path, force: true)
  end

  puts("✓ Build complete for #{cfg[:dest_os_arch]}: #{bin_binary_path}")
end

puts
puts '🏁 Done!'
