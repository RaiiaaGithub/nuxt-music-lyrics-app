<template>
    <main>
        <section class="hero">
            <div class="card">
                <div class="card-header">
                    <h2>LyriLens</h2>
                    <p>Never forget a chord with the <b>LyriLens</b> songbook builder</p>
                </div>
                <div class="card-content">
                    <input
                        ref="input"
                        v-model="inputModel"
                        type="text"
                        name="q"
                        placeholder="When I was your man - Bruno Mars"
                        :disabled="loading"
                        @keypress.enter="submit"
                    />
                    <button type="button" @click="submit" :disabled="loading">
                        <!-- Place a loading icon here -->
                        Search
                    </button>
                </div>
            </div>
        </section>
    </main>
</template>

<script setup lang="ts">
import { useSongStore } from '@/app/shared/stores/song';
import type { Song } from '@/app/shared/types/song';
import { safe } from '@/app/shared/utils/async';
import axios from 'axios';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const store = useSongStore();
const router = useRouter();

const loading = ref<boolean>(false);

const input = ref<HTMLInputElement | undefined>(undefined);
const inputModel = defineModel<string>();

onMounted(() => {
    input.value?.focus();
});

async function submit() {
    if (loading.value) {
        return;
    }
    const [res, err] = await safe(async () => {
        return await axios.get<Song>('http://localhost:8080/api/lyrics', {
            params: { q: inputModel.value }
        });
    });

    if (err) {
        console.error(err);
        loading.value = false;
        return;
    }

    store.song = res?.data;
    router.push('/editor');
    loading.value = false;
}
</script>

<style scoped>
main {
    height: 100%;
    width: 100%;
}

.hero {
    width: 100%;
    height: 100%;

    display: grid;
    justify-content: stretch;
    align-content: center;
}

.card {
    border: 1px solid var(--color-border);
    border-radius: var(--size-base);

    padding-inline: var(--padding-large);
    padding-block: var(--padding-base);

    display: grid;
    gap: var(--gap-large);

    max-width: 600px;
    width: min(600px, 100%);
    justify-self: center;
}

.card-header h2 {
    font-size: var(--font-size-large);
}

.card-content {
    display: grid;
    grid-template-columns: 1fr max-content;

    gap: var(--gap-small);
}
</style>
