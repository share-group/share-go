package echoutil

import loggerFactory "github.com/share-group/share-go/provider/logger"

func Echo() {
	loggerFactory.GetLogger("echo-tool").Info("Echo")
}
