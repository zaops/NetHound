import { ContainerLayout } from '@/components/container-layout'
import { GoBackToTop } from '@/components/goback-to-top'
import { Identity, Location, ProAndCareer, Speciality, Tecs, Trace } from '@/features/about'
import { FavSinger } from '@/features/about/fav-singer'

export default function About() {
  return (
    <ContainerLayout>
      <div className="grid max-md:gap-4 grid-cols-24 grid-rows-3 md:grid-rows-[auto_auto_1fr] gap-6 h-full max-md:flex flex-col">
        {/* 自我介绍 */}
        <Identity />
        {/* 追求 */}
        <Trace />
        {/* 技术力 */}
        <Tecs />
        {/* 专业和生涯 */}
        <ProAndCareer />

        <div className="grid grid-flow-col grid-cols-24 grid-rows-4 gap-6 col-span-24 max-md:flex flex-col animate-fade-up animate-ease-in-out animate-delay-[600ms]">
          {/* 位置 */}
          <Location />
          {/* 特长 */}
          <Speciality />
          {/* 喜欢的歌手 */}
          <FavSinger />
        </div>
      </div>

      <GoBackToTop />
    </ContainerLayout>
  )
}
