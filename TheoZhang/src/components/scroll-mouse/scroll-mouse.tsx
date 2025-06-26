export const ScrollMouse = () => {
  return (
    <div className="relative grid h-[30px] w-[20px] justify-center rounded-full border-2 border-primary/30 dark:border-borderColor/30 pt-2 md:h-[38px] md:w-[26px]">
      <div className="h-[5px] w-[2px] animate-scroll-mouse rounded-full bg-primary/30 dark:bg-bgPrimary/30 md:h-[7px]"></div>
    </div>
  )
}
