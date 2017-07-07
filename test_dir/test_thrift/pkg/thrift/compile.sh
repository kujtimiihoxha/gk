# Install thrift macOs.
# 1. brew install boost
# 2. brew install libevent
# 3. brew install thrift
# For linux see : https://thrift.apache.org/download
thrift -r --gen "go:package_prefix=github.com/kujtimiihoxha/gk/test_dir/test_thrift/pkg/thrift/gen-go/,thrift_import=github.com/apache/thrift/lib/go/thrift" test_thrift.thrift


