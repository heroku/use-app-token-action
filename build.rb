#!/usr/bin/env ruby
# frozen_string_literal: true

require 'digest'
require 'fileutils'
require 'optparse'

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

options = { dry_run: false }

OptionParser.new do |opt|
  opt.on('-d', '--dry-run', 'Skips actually replacing binaries') { options[:dry_run] = true }
end.parse!

puts(format('🚧 Starting!%s', options[:dry_run] ? ' -- DRY RUN!!! (No binaries will be updated)' : ''))
FileUtils.rm_rf('tmp')

build_configs.each do |cfg|
  new_binary_path = "tmp/main-#{cfg[:dest_os_arch]}"
  old_binary_path = "bin/main-#{cfg[:dest_os_arch]}"
  build_cmd = "#{cfg[:target_os_arch]} go build -ldflags=\"-s -w\" -o #{new_binary_path} cmd/get-gh-app-token/main.go"

  puts("Building binary for #{cfg[:dest_os_arch]}")
  build_cmd_result = system(build_cmd)

  unless build_cmd_result
    raise "An error occurred while creating the binary for #{cfg[:dest_os_arch]}: #{new_binary_path}"
  end

  if File.exist?(old_binary_path)
    new_binary_check_sum = Digest::SHA1.file(new_binary_path).hexdigest
    old_binary_check_sum = Digest::SHA1.file(old_binary_path).hexdigest

    if new_binary_check_sum == old_binary_check_sum
      puts("  Binary for #{cfg[:dest_os_arch]} is up to date. 👍")
    else
      puts(
        "  Replacing binary for #{cfg[:dest_os_arch]} with newer version",
        "  SHA256: (old) #{old_binary_check_sum} --> ✅ (new) #{new_binary_check_sum}"
      )
      FileUtils.mv(new_binary_path, old_binary_path, force: true) unless options[:dry_run]
    end
  else
    puts("  Adding binary for #{cfg[:dest_os_arch]}")
    FileUtils.mv(new_binary_path, old_binary_path, force: true) unless options[:dry_run]
  end

  puts("✓ Build complete for #{cfg[:dest_os_arch]}: #{old_binary_path}")
end

puts('', format('🏁 Done!%s', options[:dry_run] ? ' -- DRY RUN!!! (No binaries were updated)' : ''))
