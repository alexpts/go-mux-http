import { defineClientConfig } from '@vuepress/client'

const autoDetectLocale = () => navigator.language.indexOf("ru") !== -1 ? '/ru/' : '/en/'

export default defineClientConfig({
    enhance({ router }) {
        if (!__VUEPRESS_SSR__) {
            router.addRoute(
                {
                    name: "auto redirect to locale main page",
                    path: '',
                    redirect: autoDetectLocale()
                },
            )
        }
    },
})
