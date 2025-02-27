<script setup lang="ts">
import { onMounted, reactive, ref, useTemplateRef } from 'vue';
import Editor from './editor/index.vue'
import MemosCard from './memos-card/index.vue'
import {
    CalendarDaysIcon,
    CalendarIcon,
    ChevronLeftIcon,
    ChevronRightIcon,
    ListTodoIcon,
    SearchIcon
} from 'lucide-vue-next';
import dayjs from 'dayjs';
import { queryMemosByFilter, queryTags } from '../api/memo';
import { RequestCode } from '../api';
import Switch from '../components/Switch.vue';

const memosList = ref<({
    id: string
    author: string
    updateDate: string
    content: string
})[]>([])
const isLoading = ref(false)
const isMore = ref(true)
const tags = ref<{ id: string, name: string }[]>([])
const queryFilter = reactive<{
    pageNo: number
    pageSize: number
    tags: string[]
}>({
    pageNo: 1,
    pageSize: 10,
    tags: []
})

const dateList = reactive({
    current: {
        year: 2025,
        month: 1,
        day: 12
    },
    list: [
        [29, 30, 31, 1, 2, 3, 4,],
        [5, 6, 7, 8, 9, 10, 11,],
        [12, 13, 14, 15, 16, 17, 18,],
        [19, 20, 21, 22, 23, 24, 25,],
        [26, 27, 28, 29, 30, 31, 1]
    ]
})

async function onViewVisible(entries: IntersectionObserverEntry[]) {
    if (!entries.length) {
        return
    }
    const entry = entries[0]
    if (!entry.isIntersecting) {
        return
    }

    const _data = await queryMemosByFilter(queryFilter)
    const data = dataTransform(_data.data)
    if (data.length < 10) {
        isMore.value = false
        observer.unobserve(loadingEle.value!)
    }

    queryFilter.pageNo += 1
    memosList.value = [...memosList.value, ...data]
}

const loadingEle = useTemplateRef("loading")
// 监听底部加载
const observer = new IntersectionObserver(onViewVisible)

function dataTransform(data: any) {
    return data.map((item: any) => (
        { id: item.id, author: item.username, updateDate: dayjs(item.updateAt).format("YYYY/MM/DD HH:mm"), content: item.content }
    ))
}

function onCreate(data: any) {
    // console.log(data)
    memosList.value = [...dataTransform([data]), ...memosList.value]
}

function onUpdate(memo: any) {
    memosList.value = memosList.value.map((_memo) => {
        if (_memo.id != memo.id) {
            return _memo
        } else {
            _memo.content = memo.content
            return _memo
        }
    })
}

async function initMemos() {
    memosList.value = []
    isLoading.value = true
    isMore.value = true
    const _data = await queryMemosByFilter(queryFilter)

    if (_data.code === RequestCode.REQUEST_ERROR) {
        isLoading.value = false
        isMore.value = false
        return
    }

    memosList.value = dataTransform(_data.data)
    if (memosList.value.length < 10) {
        isMore.value = false
    }
    isLoading.value = false

    queryFilter.pageNo += 1

    if (isMore.value) {
        observer.observe(loadingEle.value!)
    }
}

async function initTags() {
    const { code, data } = await queryTags()
    if (code === RequestCode.REQUEST_ERROR) {
        return
    }
    tags.value = data
}

const onTagClicked = async (tagId: string) => {
    if (queryFilter.tags.includes(tagId)) {
        queryFilter.tags = queryFilter.tags.filter(_tagId => _tagId != tagId)
    } else {
        queryFilter.tags.push(tagId)
    }
    queryFilter.pageNo = 1
    queryFilter.pageSize = 10
    initMemos()
}

onMounted(async () => {
    initMemos()
    initTags()
})
</script>

<template>
    <!-- 首页 -->
    <div class="max-w-896px mx-auto flex gap-16px">
        <div class="grow">
            <Editor mode="create" @create="onCreate" />
            <div class="mt-16px">
                <!-- <div v-if="isLoading"
                    class="mx-auto relative animate-spin w-20px h-20px rounded-1/2 border-2px bg-gradient-conic bg-gradient-from-blue bg-gradient-to-lime before:(content-[''] absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rounded-1/2 w-75% h-75% bg-white)">
                </div> -->
                <!-- memos卡片列表 -->
                <MemosCard v-for="memo in memosList" @update="onUpdate" :key="memo.id" v-bind="memo"
                    class="first:mt-0 mt-8px" />

                <div v-if="isMore" ref="loading" class="flex flex-col items-center mt-16px">
                    <div
                        class="mx-auto relative animate-spin w-20px h-20px rounded-1/2 border-2px bg-gradient-conic bg-gradient-from-blue bg-gradient-to-lime before:(content-[''] absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rounded-1/2 w-75% h-75% bg-white)">
                    </div>
                    <div class="text-#aaaaaa">加载中...</div>
                </div>
                <div v-else class="flex flex-col items-center mt-16px">
                    <div class="text-#aaaaaa">没有更多了，加紧创作吧</div>
                </div>
            </div>
        </div>
        <!-- 侧边面板 -->
        <div class="w-200px shrink-0">
            <div class="sticky top-0">
                <!-- 时间面板 -->
                <div class="w-full p-8px bg-#fcfcfa rounded-8px border-(1px #ececec) shadow-[0_0_2px_#ececec]">
                    <div class="flex items-center">
                        <div class="text-#808080 cursor-pointer flex items-center">
                            <CalendarDaysIcon :size="20" />
                            <div class="text-1em ml-2px">{{ dateList.current.year }}年{{ dateList.current.month }}月</div>
                        </div>
                        <div class="ml-auto flex items-center text-#808080">
                            <ChevronLeftIcon :size="20"
                                class="cursor-pointer rounded-1/2 hover:(bg-#ffffff border-1px)" />
                            <ChevronRightIcon :size="20"
                                class="cursor-pointer rounded-1/2 hover:(bg-#ffffff border-1px)" />
                        </div>
                    </div>
                    <div class="mt-8px flex  justify-between gap-row-8px text-0.75em text-#808080">
                        <div
                            class="w-20px h-20px flex items-center justify-center cursor-pointer rounded-1/2 hover:(border-(1px #f1f1f1))">
                            日
                        </div>
                        <div
                            class="w-20px h-20px flex items-center justify-center cursor-pointer rounded-1/2 hover:(border-(1px #f1f1f1))">
                            一
                        </div>
                        <div
                            class="w-20px h-20px flex items-center justify-center cursor-pointer rounded-1/2 hover:(border-(1px #f1f1f1))">
                            二
                        </div>
                        <div
                            class="w-20px h-20px flex items-center justify-center cursor-pointer rounded-1/2 hover:(border-(1px #f1f1f1))">
                            三
                        </div>
                        <div
                            class="w-20px h-20px flex items-center justify-center cursor-pointer rounded-1/2 hover:(border-(1px #f1f1f1))">
                            四
                        </div>
                        <div
                            class="w-20px h-20px flex items-center justify-center cursor-pointer rounded-1/2 hover:(border-(1px #f1f1f1))">
                            五
                        </div>
                        <div
                            class="w-20px h-20px flex items-center justify-center cursor-pointer rounded-1/2 hover:(border-(1px #f1f1f1))">
                            六
                        </div>
                    </div>

                    <div v-for="days in dateList.list" :key="days[0]"
                        class="mt-4px flex  justify-between gap-row-8px text-0.75em text-#808080">
                        <div v-for="day in days" :key="day"
                            class="w-20px h-20px flex items-center justify-center cursor-pointer rounded-1/2 hover:(border-(1px #f1f1f1))"
                            :class="{ 'bg-#52aaa0 text-white font-600': dateList.current.day === day }">
                            {{ day }}
                        </div>
                    </div>

                    <div class="mt-8px text-#808080 text-0.75em">
                        6天记录3条
                    </div>

                    <div class="mt-8px flex gap-4px">
                        <div
                            class="flex items-center gap-4px py-2px px-8px bg-#ffffff text-0.75em border-(1px #f0f0f0) rounded-1em cursor-pointer ">
                            <ListTodoIcon :size="12" />
                            <span>待办</span>
                            <span>3</span>
                        </div>
                        <div
                            class="flex items-center gap-4px py-2px px-8px bg-#ffffff text-0.75em border-(1px #f0f0f0) rounded-1em cursor-pointer ">
                            <CalendarIcon :size="12" />
                            <span>今日</span>
                            <span>2</span>
                        </div>
                    </div>
                </div>
                <!-- 搜索面板 -->
                <div class="mt-16px w-full">
                    <div class="flex items-center bg-white h-2em px-0.5em rounded-8px overflow-hidden">
                        <SearchIcon class="text-#808080" :size="20" />
                        <input class="outline-none w-full h-full ml-4px" placeholder="搜索备忘录" />
                    </div>
                    <!-- 标签区域 -->
                    <div class="mt-16px">
                        <div class="flex items-center">
                            <span class="text-#808080">标签集</span>
                            <div class="ml-auto">
                                <Switch />
                            </div>
                        </div>
                        <div class="mt-8px flex gap-8px gap-row-2px flex-wrap">
                            <template v-for="tag in tags" :key="tag.id">
                                <div :data-tagId="tag.id" @click="onTagClicked(tag.id)"
                                    class="flex items-center cursor-pointer border-1px rounded-0.75em pl-0.5em pr-0.5em">
                                    <span class="text-1em text-#808080">#</span>
                                    <span class="ml-2px text-0.8em">{{ tag.name }}</span>
                                    <span v-if="queryFilter.tags.includes(tag.id)"
                                        class="i-lucide:x text-#acacac ml-0.3em"></span>
                                </div>
                            </template>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>