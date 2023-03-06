package pac

import (
	"fmt"
	"github.com/darki73/pac-manager/pkg/logger"
	"strings"
	"text/template"
)

var pacTemplate = `globalThis.fallback = 'DIRECT';

globalThis.proxies = [
	{{ .ProxiesList }}
];

globalThis.rules = {
	{{ .RulesList }}
};

function GetRandomProxy(proxies) {
    return proxies[Math.floor(Math.random() * proxies.length)]
}

function FindProxyForURL(url, host) {
    let level = 0;
    let position = 0;
    
    const proxy = GetRandomProxy(proxies);
    let noWildCardRules = [];
    
    Object.keys(rules).forEach((rule) => {
        const shouldProxy = rules[rule];
        if (shouldProxy) {
            noWildCardRules.push(rule.replace('*.', ''));
        }
    });
    
    while (position >= 0) {
        if (rules.hasOwnProperty(host)) {
            if (rules[host]) {
                return proxy;
            }
            return fallback;
        }
        
        if (noWildCardRules.includes(host)) {
            return proxy;
        }
        
        if (position >= 0 && level === 1) {
            if (rules['*.' + host]) {
                return proxy;
            }
        }
        
        position = host.indexOf('.');
        level += 1;
        host = host.slice(position + 1);
    }

	return fallback;
}
`

// Generate generates the pac file.
func (pac *Pac) Generate(proxies []string, rules map[string]bool) {
	tpl, err := template.New("").Parse(pacTemplate)
	if err != nil {
		logger.Panicf("pac:generate", "failed to parse pac template - %s", err.Error())
	}

	processedProxies := make([]string, 0, len(proxies))

	for _, proxy := range proxies {
		processedProxies = append(processedProxies, fmt.Sprintf("'%s',", proxy))
	}

	processedRules := make([]string, 0, len(rules))

	for rule, proxy := range rules {
		processedRules = append(processedRules, fmt.Sprintf("'%s': %t,", rule, proxy))
	}

	content := struct {
		ProxiesList string
		RulesList   string
	}{
		ProxiesList: strings.Join(processedProxies, "\n"),
		RulesList:   strings.Join(processedRules, "\n"),
	}

	file, err := pac.GetStorage().OpenFile(pac.GetName())
	if err != nil {
		logger.Fatalf("pac:generate", "failed to open pac file - %s", err.Error())
	}
	defer file.Close()

	if err = tpl.Execute(file, content); err != nil {
		logger.Fatalf("pac:generate", "failed to generate pac output - %s", err.Error())
	}
}
