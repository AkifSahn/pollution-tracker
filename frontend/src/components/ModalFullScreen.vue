<template>
    <transition name="modal">
        <div v-if="show" class="modal-backdrop">
            <div class="modal-container">
                <div class="modal-header">
                    <h3 class="modal-title">{{ title }}</h3>
                    <button @click="$emit('close')" class="modal-close">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                            stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
                <div class="modal-content">
                    <slot></slot>
                </div>
            </div>
        </div>
    </transition>
</template>

<script>
export default {
    name: 'ModalFullScreen',
    props: {
        show: {
            type: Boolean,
            default: false
        },
        title: {
            type: String,
            default: ''
        }
    },
    watch: {
        show(newVal) {
            if (newVal) {
                document.body.style.overflow = 'hidden';
            } else {
                document.body.style.overflow = '';
            }
        }
    },
    beforeUnmount() {
        document.body.style.overflow = '';
    }
}
</script>

<style scoped>
.modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 100;
    display: flex;
    align-items: center;
    justify-content: center;
}

.modal-container {
    position: relative;
    width: 90%;
    height: 90%;
    max-width: 1400px;
    background-color: var(--card-bg);
    border-radius: 8px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.25);
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 24px;
    background-color: var(--primary-color);
    color: var(--header-text);
    border-bottom: 1px solid var(--card-shadow);
}

.modal-title {
    font-weight: 600;
    font-size: 1.2rem;
    margin: 0;
}

.modal-close {
    background: transparent;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--header-text);
    padding: 4px;
    border-radius: 4px;
    transition: background-color 0.2s;
}

.modal-close:hover {
    background-color: rgba(255, 255, 255, 0.2);
}

.modal-content {
    flex-grow: 1;
    padding: 0;
    overflow: auto;
    height: calc(100% - 60px);
}

/* Transitions */
.modal-enter-active,
.modal-leave-active {
    transition: opacity 0.3s;
}

.modal-enter-from,
.modal-leave-to {
    opacity: 0;
}
</style>

