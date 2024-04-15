<script setup>
  const task = defineModel();
</script>

<script>
  export default {
    emits: ['show-task', 'delete-task'],
    methods: {
      goToTask() {
        this.$emit('show-task');
      },
      deleteTask() {
        this.$emit('delete-task')
      }
    }
  }
</script>

<template>
    <div class="row-cols-lg-6">
      <div class="col">
        <button class="blk-btn selector" @click="goToTask">Перейти к задаче</button>
      </div>
      <div class="col">
        <p class="annot">Название:</p>
        <p class="value">{{ task.title }}</p>
      </div>
      <div class="col">
        <p class="annot">Статус:</p>
        <p class="value">{{ task.status ? 'Завершено' : 'Черновик' }}</p>
      </div>
      <div class="col">
        <p class="annot">Последние изменения:</p>
        <p class="value">{{ new Date(task.last_change).toDateString() }}</p>
      </div>
      <div class="col">
        <p class="annot">Тип задачи:</p>
        <p class="value">{{ task.taskType === 'individual' ? 'Индивидуальная' : 'Групповая' }}</p>
      </div>
      <div class="col">
        <button class="deleter">
          <img alt="delete" class="delete" @click="deleteTask" src="/delete.png">
        </button>
      </div>
    </div>
</template>

<style scoped>
  @import "../../style.css";
  @import "../../assets/personal.css";
  @import "../../assets/tasklist.css";
</style>