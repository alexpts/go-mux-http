// https://v2.vuepress.vuejs.org/reference/default-theme/config.html

import { defaultTheme } from '@vuepress/theme-default'
import { searchPlugin } from '@vuepress/plugin-search'
import { nprogressPlugin } from '@vuepress/plugin-nprogress'
import umlPlugin from 'markdown-it-plantuml'
import { viteBundler } from '@vuepress/bundler-vite'
import { defineUserConfig } from 'vuepress'

let ruThemeConfig = {
    selectLanguageText: 'Языки',
    selectLanguageName: 'Русский',
    editLinkText: 'Редактировать',
    contributorsText: 'Авторы',
    lastUpdatedText: 'Обновлено',

    navbar: [
        { text: 'Главная', link: '/ru/' },
        { text: 'Быстрый старт', link: '/ru/guide/' },
    ],
};

let enThemeConfig = {
    selectLanguageText: 'Languages',
    selectLanguageName: 'English',

    navbar: [
        { text: 'Home', link: '/en/' },
        { text: 'Quick Start', link: '/en/guide/' },
    ],
};


module.exports = defineUserConfig({
    base: '/go-mux-http/', // github pages sub-url

    plugins: [
        searchPlugin(),
        nprogressPlugin()
    ],

    bundler: viteBundler({
        viteOptions: {},
        vuePluginOptions: {},
    }),

    locales: {
        '/ru/': {
            lang: 'ru',
        },
        '/en/': {
            lang: 'en',
        },
    },

    theme: defaultTheme({
        search: true,

        locales: {
            '/': ruThemeConfig,
            '/ru/': ruThemeConfig,
            '/en/': enThemeConfig,
        },

        repo: 'alexpts/go-mux-http',

        //docsRepo: 'alexpts/go-next-docs/',
        docsBranch: 'master',
        docsDir: 'docs/src',
    }),

    extendsMarkdown: (md) => {
        md.use(umlPlugin);   // required by PalmUML
    },
})
