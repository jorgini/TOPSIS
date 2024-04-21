<script setup>
  const task = defineModel();
  const emits = defineEmits(["show-task", "delete-task", "show-final"]);
</script>

<script>
  export default {
    emits: ['show-task', 'delete-task', 'show-final'],
    methods: {
      goToTask() {
        this.$emit('show-task');
      },
      goToFinal() {
        this.$emit('show-final');
      },
      deleteTask() {
        this.$emit('delete-task')
      }
    }
  }
</script>

<template>
    <div class="row-cols-6">
      <div class="col-3">
        <div class="btns">
          <button class="blk-btn" @click="goToTask">Перейти к задаче</button>
          <button v-if="task.status === true" class="cl-btn" @click="goToFinal">Перейти к отчету</button>
        </div>
      </div>
      <div class="col-2">
        <p class="annot">Название:</p>
        <p class="value">{{ task.title }}</p>
      </div>
      <div class="col-2">
        <p class="annot">Статус:</p>
        <p class="value">{{ task.status ? 'Завершено' : 'Черновик' }}</p>
      </div>
      <div class="col-2">
        <p class="annot">Последние изменения:</p>
        <p class="value">{{ new Date(task.last_change).toDateString() }}</p>
      </div>
      <div class="col-2">
        <p class="annot">Тип задачи:</p>
        <p class="value">{{ task.taskType === 'individual' ? 'Индивидуальная' : 'Групповая' }}</p>
      </div>
      <div class="col-1">
        <button class="deleter">
          <img alt="delete" class="delete" @click="deleteTask" src="/delete.png">
        </button>
      </div>
    </div>
</template>

<style scoped>
  @import "../../style.css";
  @import "../../assets/tasklist.css";
  .col-3 {
    display: flex;
    justify-content: center;
  }

  .col-3 > button {
    margin: auto;
  }

  .col-1 {
    display: flex;
    justify-content: center;
  }

  .col-1 > button {
    margin: auto;
  }


</style>