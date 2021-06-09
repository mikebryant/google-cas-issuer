package framework

import (
	cmversioned "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	"github.com/jetstack/google-cas-issuer/test/e2e/framework/config"
	"github.com/jetstack/google-cas-issuer/test/e2e/framework/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
)

type Framework struct {
	BaseName string

	KubeClientSet    kubernetes.Interface
	CMClientSet      cmversioned.Interface
	DynamicClientSet dynamic.Interface

	DiscoveryClient discovery.DiscoveryInterface
	Mapper          *restmapper.DeferredDiscoveryRESTMapper

	config *config.Config
	helper *helper.Helper
}

func NewDefaultFramework(baseName string) *Framework {
	return NewFramework(baseName, config.GetConfig())
}

func NewFramework(baseName string, config *config.Config) *Framework {
	f := &Framework{
		BaseName: baseName,
		config:   config,
	}

	JustBeforeEach(f.BeforeEach)

	return f
}

func (f *Framework) BeforeEach() {
	By("Creating a kubernetes client")
	clientConfigFlags := genericclioptions.NewConfigFlags(true)
	clientConfigFlags.KubeConfig = &f.config.KubeConfigPath
	config, err := clientConfigFlags.ToRESTConfig()
	Expect(err).NotTo(HaveOccurred())

	f.KubeClientSet, err = kubernetes.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())

	By("Creating a cert-manager client")
	f.CMClientSet, err = cmversioned.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())

	By("Creating a dynamic client")
	f.DynamicClientSet, err = dynamic.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())

	By("Creating a discovery client")
	f.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(config)
	Expect(err).NotTo(HaveOccurred())

	f.Mapper = restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(f.DiscoveryClient))

	f.helper = helper.NewHelper(f.CMClientSet, f.KubeClientSet, f.DynamicClientSet)
}

func (f *Framework) Helper() *helper.Helper {
	return f.helper
}

func (f *Framework) Config() *config.Config {
	return f.config
}

func CasesDescribe(text string, body func()) bool {
	return Describe("[jetstack google-cas-issuer] "+text, body)
}
