package registrar

import "github.com/google/wire"

// ProviderSet is Registry providers.
var ProviderSet = wire.NewSet(NewEtcdRegistrar)
