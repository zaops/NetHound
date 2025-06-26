'use client'

import { GoBackToTop } from '@/components/goback-to-top'
import IconSelf from '@/components/icons/icon-self'
import { ScrollMouse } from '@/components/scroll-mouse'
import KlButton from '@/components/ui/button'
import {
  GITHUB_PERSONAL_PAGE_LINK,
  HOME_NAME_TEXT,
  HOME_PRE_TEXT,
  HOME_SLOGAN_TEXT,
  HOME_TPYE_EN_TEXT,
  HOME_TPYE_ZH_TEXT,
  TWITTER_PERSONAL_PAGE_LINK
} from '@/constants/info'
import { PATHS } from '@/constants/path'
import { clm } from '@/utils'
import Link from 'next/link'
import React, { useEffect } from 'react'
import { TypeAnimation } from 'react-type-animation'

const Section = () => {
  const colorList = {
    vueTheme: 'from-[#43D392] to-[#647EFF]',
    reactTheme: 'from-[#58C4DC] to-[#3178C6]'
  }

  const [theme, setTheme] = React.useState<string>(colorList.vueTheme)

  useEffect(() => {
    function getRandomTheme() {
      const randomCount = Math.random()
      const theRandomTheme = randomCount < 0.5 ? colorList.vueTheme : colorList.reactTheme
      // console.log(randomCount)
      return theRandomTheme
    }

    setTheme(getRandomTheme())
  }, [colorList.vueTheme, colorList.reactTheme])

  return (
    <>
      <div className="flex flex-col gap-4 px-4 mx-auto text-primary dark:text-darkprimary tracking-wider">
        {/* 首页文字 */}
        <>
          <span className="text-5xl max-sm:text-2xl animate-fade-up animate-ease-in-out">
            {HOME_PRE_TEXT}
          </span>
          <span
            className={clm(
              'text-8xl max-sm:text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r animate-fade-up animate-ease-in-out animate-delay-[200ms]',
              theme
            )}
          >
            {HOME_NAME_TEXT}
          </span>

          <TypeAnimation
            sequence={[HOME_TPYE_ZH_TEXT, 1000, HOME_TPYE_EN_TEXT, 2000, () => {}]}
            wrapper="span"
            cursor={true}
            repeat={Infinity}
            speed={8}
            className="text-5xl max-sm:text-2xl flex break-all animate-fade-up animate-ease-in-out animate-delay-[400ms]"
          />

          <span className="flex gap-2 text-5xl flex-wrap max-sm:text-2xl animate-fade-up animate-ease-in-out animate-delay-[600ms]">
            <span>喜欢</span>
            <span>
              <span className="font-bold text-[#43D392]">Linux</span>
              <span>、</span>
            </span>
            <span>
              <span className="font-bold text-[#58C4DC]">Python</span>
              <span>、</span>
            </span>
            <span>
              <span className="font-bold text-[#3178C6] mr-4">Workflow</span>
              <span>🙈 ~</span>
            </span>
          </span>

          <span className="text-2xl max-sm:text-lg text-secondary dark:text-darksecondary animate-fade-up animate-ease-in-out animate-delay-[800ms]">
            {HOME_SLOGAN_TEXT}
          </span>
        </>

        {/* 底部按钮 */}
        <div className="flex gap-4 text-sm font-medium animate-fade-up animate-ease-in-out animate-delay-[1000ms]">
          <Link href={PATHS.SITE_NOTE}>
            <KlButton>我的文章</KlButton>
          </Link>

          <Link href={PATHS.SITE_ABOUT}>
            <KlButton>关于我</KlButton>
          </Link>

          <Link href={GITHUB_PERSONAL_PAGE_LINK} target="_blank" className="flex">
            <KlButton isIconOnly={true} className="h-full">
              <IconSelf iconName="icon-[fa6-brands--github]" />
            </KlButton>
          </Link>

          <Link href={TWITTER_PERSONAL_PAGE_LINK} target="_blank" className="flex">
            <KlButton isIconOnly={true} className="h-full">
              <IconSelf iconName="icon-[fa6-brands--x-twitter]" />
            </KlButton>
          </Link>
        </div>

        <div className="absolute inset-x-0 bottom-8 grid place-content-center md:bottom-12 animate-fade-up animate-ease-in-out animate-delay-[1200ms]">
          <ScrollMouse />
        </div>
      </div>

      <GoBackToTop />
    </>
  )
}

export default Section
