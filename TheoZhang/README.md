# <div align="center">🚀 TheoZhang Blog网站</div>

<div align="center">
  <img src="./public/TheoZhang.svg" alt="TheoZhang Logo" width="200" height="200" />
</div>

<div align="center">

  [![GitHub](https://img.shields.io/badge/GitHub-TheoZhang-blue?logo=github)](https://github.com/zaops)
  [![Next.js](https://img.shields.io/badge/Next.js-15.3.2-black?logo=next.js)](https://nextjs.org/)
  [![React](https://img.shields.io/badge/React-19.0.0-blue?logo=react)](https://reactjs.org/)
  [![TypeScript](https://img.shields.io/badge/TypeScript-5.0-blue?logo=typescript)](https://www.typescriptlang.org/)
  [![Tailwind CSS](https://img.shields.io/badge/Tailwind%20CSS-4.0-38B2AC?logo=tailwind-css)](https://tailwindcss.com/)
</div>

---

## 📖 项目简介

这是 TheoZhang 的个人博客网站，一个现代化的全栈 Web 应用程序。网站包含个人介绍、技能展示、博客文章管理等功能，旨在记录个人成长历程和技术分享。

### ✨ 主要特色

- 🎨 **现代化设计** - 使用 Tailwind CSS 和 HeroUI 组件库打造精美界面
- 🌙 **深色模式** - 支持明暗主题切换，提供更好的用户体验
- 📱 **响应式设计** - 完美适配桌面端和移动端设备
- ✍️ **Markdown 编辑器** - 集成 ByteMD 编辑器，支持实时预览
- 🔍 **代码高亮** - 使用 highlight.js 提供代码语法高亮
- 🎭 **动画效果** - 使用 Framer Motion 和 tailwindcss-animated 实现流畅动画
- 🏷️ **标签管理** - 支持文章标签分类和管理
- 🔧 **后台管理** - 提供完整的内容管理后台

---

## 🛠️ 技术栈

### 前端技术
- **框架**: Next.js 15.3.2 (App Router)
- **UI 库**: React 19.0.0
- **语言**: TypeScript 5.0
- **样式**: Tailwind CSS 4.0
- **组件库**: HeroUI 2.8.0-beta.7
- **动画**: Framer Motion 12.17.3
  - **图标**: Iconify (多个图标包)

### 开发工具
- **包管理器**: PNPM
- **代码规范**: ESLint
- **字体**: Geist Sans & Geist Mono
- **主题**: next-themes (主题切换)

### 核心依赖
- **Markdown 编辑器**: ByteMD 1.22.0
- **代码高亮**: highlight.js 11.11.1
- **表格组件**: TanStack React Table 8.21.3
- **工具库**: lodash-es 4.17.21
- **样式工具**: clsx, tailwind-merge

---

## 📁 项目结构

```
klausjin/
├── public/                    # 静态资源文件
│   ├── about/                # 关于页面相关图片
│   ├── klausjin-dark.svg    # 深色模式 Logo
│   └── klausjin-light.svg   # 浅色模式 Logo
├── src/
│   ├── app/                  # Next.js App Router 页面
│   │   ├── (backend)/       # 后台管理页面组
│   │   │   └── admin/       # 管理员页面
│   │   ├── (site)/          # 前台网站页面组
│   │   │   ├── about/       # 关于页面
│   │   │   ├── note/        # 博客文章页面
│   │   │   └── @modal/      # 模态框路由
│   │   ├── api/             # API 路由
│   │   ├── layout.tsx       # 根布局组件
│   │   └── providers.tsx    # 全局提供者组件
│   ├── backend/             # 后台管理相关组件
│   │   ├── admin.tsx        # 管理员组件
│   │   ├── backend-home/    # 后台首页
│   │   ├── backend-note/    # 后台文章管理
│   │   └── backend-tag/     # 后台标签管理
│   ├── components/          # 可复用组件
│   │   ├── backend-avatar/  # 后台头像组件
│   │   ├── backend-navbar/  # 后台导航栏
│   │   ├── card-layout/     # 卡片布局
│   │   ├── container-layout/ # 容器布局
│   │   ├── footer/          # 页脚组件
│   │   ├── goback-to-top/   # 回到顶部按钮
│   │   ├── icons/           # 图标组件
│   │   ├── logo/            # Logo 组件
│   │   ├── markdown/        # Markdown 相关组件
│   │   ├── navbar/          # 导航栏组件
│   │   ├── note-card/       # 文章卡片
│   │   ├── providers/       # 提供者组件
│   │   ├── rolling-text/    # 滚动文字效果
│   │   ├── scroll-mouse/    # 滚动鼠标动画
│   │   ├── toggleMode/      # 主题切换按钮
│   │   └── ui/              # 基础 UI 组件
│   ├── constants/           # 常量定义
│   │   ├── assets.ts        # 资源常量
│   │   ├── backend-info.ts  # 后台信息常量
│   │   ├── backend-path.ts  # 后台路径常量
│   │   ├── info.ts          # 网站信息常量
│   │   └── path.ts          # 路径常量
│   ├── features/            # 功能模块
│   │   ├── about/           # 关于页面功能
│   │   ├── home/            # 首页功能
│   │   └── note/            # 博客功能
│   ├── hooks/               # 自定义 Hooks
│   │   ├── useForceUpdate.ts # 强制更新 Hook
│   │   ├── useScroll.ts     # 滚动 Hook
│   │   └── useToast.ts      # 提示 Hook
│   ├── lib/                 # 工具库
│   │   └── hero.ts          # HeroUI 配置
│   ├── md/                  # Markdown 文件
│   │   └── test.md          # 测试文章
│   ├── styles/              # 样式文件
│   │   ├── globals.css      # 全局样式
│   │   ├── custom.css       # 自定义样式
│   │   ├── custom-dark.css  # 深色模式样式
│   │   ├── github.css       # GitHub 样式
│   │   └── github-dark.css  # GitHub 深色样式
│   └── utils/               # 工具函数
│       ├── normal.ts        # 通用工具
│       ├── style.ts         # 样式工具
│       └── note/            # 文章相关工具
├── eslint.config.mjs        # ESLint 配置
├── next.config.ts           # Next.js 配置
├── package.json             # 项目依赖
├── postcss.config.mjs       # PostCSS 配置
├── tsconfig.json            # TypeScript 配置
└── pnpm-lock.yaml          # PNPM 锁定文件
```

---

## 🚀 快速开始

### 环境要求

- Node.js 18.0 或更高版本
- PNPM 8.0 或更高版本（推荐）

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/KlausJins/klausjin.git
   cd klausjin
   ```

2. **安装依赖**
   ```bash
   pnpm install
   ```

3. **启动开发服务器**
   ```bash
   pnpm dev
   ```

4. **打开浏览器**

   访问 [http://localhost:3000](http://localhost:3000) 查看网站

### 可用脚本

```bash
# 开发模式（使用 Turbopack）
pnpm dev

# 构建生产版本
pnpm build

# 启动生产服务器
pnpm start

# 代码检查
pnpm lint
```

---

## 📄 页面功能

### 🏠 首页 (/)
- 个人介绍和欢迎信息
- 打字机效果的职业描述
- 响应式设计，适配各种设备

### 👨‍💻 关于页面 (/about)
- **个人身份**: 详细的自我介绍
- **技能展示**: 前端、后端、设计工具等技能图标
- **追求理念**: 滚动文字展示个人追求
- **专业背景**: 教育和职业经历
- **个人特色**: 编程特长和兴趣爱好
- **地理位置**: 当前所在城市
- **喜爱歌手**: 个人音乐品味

### 📝 博客页面 (/note)
- 文章列表展示
- 标签筛选功能
- 文章搜索功能
- 模态框文章详情
- 响应式文章卡片

### 🔧 后台管理 (/admin)
- **文章管理**: 创建、编辑、删除文章
- **标签管理**: 管理文章标签
- **Markdown 编辑器**: 实时预览功能
- **移动端限制**: 后台仅支持桌面端访问

---

## 🎨 设计特色

### 主题系统
- **浅色主题**: 简洁明亮的设计风格
- **深色主题**: 护眼的深色配色方案
- **自动切换**: 跟随系统主题设置

### 动画效果
- **页面过渡**: 流畅的页面切换动画
- **元素动画**: 淡入、滑入等进场动画
- **交互反馈**: 按钮悬停、点击效果
- **滚动动画**: 基于滚动位置的动画触发

### 响应式设计
- **移动优先**: 优先考虑移动端体验
- **断点适配**: 完美适配各种屏幕尺寸
- **触摸友好**: 移动端触摸操作优化

---

## 🔧 配置说明

### 网站信息配置
在 `src/constants/info.ts` 中可以修改：
- 个人昵称和网站标题
- 社交媒体链接
- 技能图标和描述
- 个人介绍文案

### 主题配置
在 `src/styles/globals.css` 中可以自定义：
- 颜色变量
- 字体设置
- 动画效果
- 响应式断点

### Next.js 配置
在 `next.config.ts` 中可以调整：
- 构建优化选项
- 图片优化设置
- 路由配置

---

## 📦 部署指南

### Vercel 部署（推荐）

1. 将代码推送到 GitHub
2. 在 [Vercel](https://vercel.com) 中导入项目
3. 自动部署完成

### 其他平台部署

1. **构建项目**
   ```bash
   pnpm build
   ```

2. **启动生产服务器**
   ```bash
   pnpm start
   ```

3. **静态导出**（如需要）
   ```bash
   # 在 next.config.ts 中添加 output: 'export'
   pnpm build
   ```

---

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

### 开发流程

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 代码规范

- 使用 TypeScript 进行类型检查
- 遵循 ESLint 代码规范
- 使用 Prettier 格式化代码
- 组件和函数需要添加注释

---

## 📝 更新日志

### v0.1.0 (当前版本)
- ✨ 初始版本发布
- 🎨 完成基础页面设计
- 📝 集成 Markdown 编辑器
- 🌙 支持深色模式
- 📱 响应式设计完成
- 🔧 后台管理功能

---

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

---

## 👨‍💻 作者

**TheoZhang** - 运维开发工程师

- GitHub: [@zaops](https://github.com/KlausJins)
- 项目链接: [https://github.com/zaops/TheoZhang-Blog](https://github.com/zaops/TheoZhang-Blog)

---

## 🙏 致谢

感谢以下开源项目的支持：

- [Next.js](https://nextjs.org/) - React 全栈框架
- [Tailwind CSS](https://tailwindcss.com/) - 实用优先的 CSS 框架
- [HeroUI](https://heroui.com/) - 现代化 React 组件库
- [ByteMD](https://bytemd.js.org/) - Markdown 编辑器
- [Framer Motion](https://www.framer.com/motion/) - React 动画库
- [Iconify](https://iconify.design/) - 图标解决方案

---

<div align="center">
  <p>⭐ 如果这个项目对你有帮助，请给它一个星标！</p>
  <p>💬 有任何问题或建议，欢迎提交 Issue</p>
</div>
