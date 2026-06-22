package model

import "github.com/smtdfc/nagare/core/config"

func BuildFromConfig(conf *config.Config) (bool, ChatModel) {

	currentProvider, ok := conf.Providers[conf.CurrentProvider]

	if !ok {
		return false, nil
	}

	if currentProvider.Compatible == config.OPEN_AI_COMPATIBLE {
		return true, NewOpenAICompatibleClient(&ChatModelConfig{
			BaseURL: currentProvider.BaseURL,
			APIKey:  currentProvider.APIKey,
			Model:   currentProvider.ModelName,
		})
	}

	return false, nil
}
